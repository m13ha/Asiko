package ahasend

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

var sharedPublisher *Publisher

type Publisher struct {
	client      *Client
	config      Config
	queue       chan MessageEvent
	workerCount int64
	shutdown    chan struct{}
	mu          sync.RWMutex
}

type MessageEvent struct {
	Message MessageRequest
	Retries int
}

func NewPublisher(config Config) (*Publisher, error) {
	if shared := getSharedPublisher(); shared != nil {
		return shared, nil
	}
	return newPublisher(config)
}

func newPublisher(config Config) (*Publisher, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if !config.Enabled {
		p := &Publisher{config: config}
		setSharedPublisher(p)
		return p, nil
	}

	client := NewClient(config)
	p := &Publisher{
		client:   client,
		config:   config,
		queue:    make(chan MessageEvent, config.MaxQueueSize),
		shutdown: make(chan struct{}),
	}

	go p.scaleWorkers()
	setSharedPublisher(p)
	return p, nil
}

func (p *Publisher) Publish(ctx context.Context, message MessageRequest) error {
	if !p.config.Enabled || p.client == nil || p.queue == nil {
		return nil
	}

	event := MessageEvent{Message: message}

	select {
	case p.queue <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		select {
		case p.queue <- event:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *Publisher) scaleWorkers() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.adjustWorkerCount()
		case <-p.shutdown:
			return
		}
	}
}

func (p *Publisher) adjustWorkerCount() {
	queueLen := len(p.queue)
	currentWorkers := int(atomic.LoadInt64(&p.workerCount))

	var neededWorkers int
	if queueLen > 0 {
		neededWorkers = (queueLen / p.config.EventsPerWorker) + 1
	}

	if neededWorkers > p.config.MaxWorkers {
		neededWorkers = p.config.MaxWorkers
	}

	if neededWorkers > currentWorkers {
		for i := currentWorkers; i < neededWorkers; i++ {
			go p.worker()
		}
	}
}

func (p *Publisher) worker() {
	atomic.AddInt64(&p.workerCount, 1)
	defer atomic.AddInt64(&p.workerCount, -1)

	workerID := atomic.LoadInt64(&p.workerCount)
	idleTimer := time.NewTimer(p.config.WorkerIdleTimeout)
	defer idleTimer.Stop()

	log.Info().Int64("workerId", workerID).Msg("ahasend worker created")

	for {
		select {
		case event := <-p.queue:
			idleTimer.Reset(p.config.WorkerIdleTimeout)
			p.processEvent(event)
		case <-idleTimer.C:
			log.Info().Int64("workerId", workerID).Msg("ahasend worker deleted")
			return
		case <-p.shutdown:
			log.Info().Int64("workerId", workerID).Msg("ahasend worker stopped by shutdown")
			return
		}
	}
}

func (p *Publisher) processEvent(event MessageEvent) {
	var err error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		err = p.client.SendMessage(context.Background(), event.Message)
		if err == nil {
			return
		}

		log.Warn().Err(err).Msg("ahasend message failed")

		if attempt < p.config.MaxRetries {
			time.Sleep(p.config.Backoff)
		}
	}

	log.Error().Err(err).Msg("ahasend message failed after retries")
}

func (p *Publisher) Shutdown(ctx context.Context) error {
	if !p.config.Enabled {
		return nil
	}

	close(p.shutdown)

	done := make(chan struct{})
	go func() {
		for atomic.LoadInt64(&p.workerCount) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("ahasend: shutdown completed")
		return nil
	case <-ctx.Done():
		log.Warn().Msg("ahasend: shutdown timeout")
		return ctx.Err()
	}
}

func (p *Publisher) GetWorkerCount() int64 {
	return atomic.LoadInt64(&p.workerCount)
}

func (p *Publisher) GetQueueSize() int {
	if p.queue == nil {
		return 0
	}
	return len(p.queue)
}

func (p *Publisher) GetQueueCapacity() int {
	if p.queue == nil {
		return 0
	}
	return cap(p.queue)
}

func (p *Publisher) Stats() (queueLen int, workerCount int) {
	return p.GetQueueSize(), int(p.GetWorkerCount())
}

func setSharedPublisher(p *Publisher) {
	sharedPublisher = p
}

func getSharedPublisher() *Publisher {
	return sharedPublisher
}
