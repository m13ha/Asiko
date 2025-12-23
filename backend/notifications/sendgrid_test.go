package notifications

import (
	"errors"
	"testing"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type fakeClient struct {
	responses []*rest.Response
	errors    []error
	calls     int
}

func (f *fakeClient) Send(_ *mail.SGMailV3) (*rest.Response, error) {
	idx := f.calls
	f.calls++
	var resp *rest.Response
	var err error
	if idx < len(f.responses) {
		resp = f.responses[idx]
	}
	if idx < len(f.errors) {
		err = f.errors[idx]
	}
	return resp, err
}

func makeResp(status int, body string) *rest.Response {
	return &rest.Response{StatusCode: status, Body: body}
}

// noSleep replaces the package-level sleep function during tests to avoid delays

func TestSendEmail_SuccessFirstTry(t *testing.T) {
	oldSleep := sleep
	sleep = func(d time.Duration) {}
	defer func() { sleep = oldSleep }()

	fc := &fakeClient{responses: []*rest.Response{makeResp(202, "ok")}}
	svc := NewSendGridServiceWithClient(fc, "from@example.com")

	if err := svc.SendVerificationCode("to@example.com", "1234"); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if fc.calls != 1 {
		t.Fatalf("expected 1 call, got %d", fc.calls)
	}
}

func TestSendEmail_RetryThenSuccess(t *testing.T) {
	oldSleep := sleep
	sleep = func(d time.Duration) {}
	defer func() { sleep = oldSleep }()

	fc := &fakeClient{responses: []*rest.Response{makeResp(500, "fail"), makeResp(202, "ok")}}
	svc := NewSendGridServiceWithClient(fc, "from@example.com")

	if err := svc.SendVerificationCode("to@example.com", "1234"); err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	if fc.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", fc.calls)
	}
}

func TestSendEmail_ExhaustRetries(t *testing.T) {
	oldSleep := sleep
	sleep = func(d time.Duration) {}
	defer func() { sleep = oldSleep }()

	fc := &fakeClient{responses: []*rest.Response{makeResp(500, "err1"), makeResp(502, "err2"), makeResp(503, "err3")}}
	svc := NewSendGridServiceWithClient(fc, "from@example.com")

	if err := svc.SendVerificationCode("to@example.com", "1234"); err == nil {
		t.Fatalf("expected error after exhausting retries, got nil")
	}
	if fc.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", fc.calls)
	}
}

func TestSendEmail_NoRetryOn4xx(t *testing.T) {
	oldSleep := sleep
	sleep = func(d time.Duration) {}
	defer func() { sleep = oldSleep }()

	fc := &fakeClient{responses: []*rest.Response{makeResp(400, "bad request")}}
	svc := NewSendGridServiceWithClient(fc, "from@example.com")

	if err := svc.SendVerificationCode("to@example.com", "1234"); err == nil {
		t.Fatalf("expected error on 4xx, got nil")
	}
	if fc.calls != 1 {
		t.Fatalf("expected 1 call, got %d", fc.calls)
	}
}

func TestSendEmail_RetryOnNetworkError(t *testing.T) {
	oldSleep := sleep
	sleep = func(d time.Duration) {}
	defer func() { sleep = oldSleep }()

	fc := &fakeClient{responses: []*rest.Response{nil, makeResp(202, "ok")}, errors: []error{errors.New("net error"), nil}}
	svc := NewSendGridServiceWithClient(fc, "from@example.com")

	if err := svc.SendVerificationCode("to@example.com", "1234"); err != nil {
		t.Fatalf("expected success after network retry, got error: %v", err)
	}
	if fc.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", fc.calls)
	}
}
