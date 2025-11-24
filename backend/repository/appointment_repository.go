package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	apperr "github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/responses"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppointmentRepository interface {
	Create(appointment *entities.Appointment) error
	GetAppointmentsByOwnerIDQuery(ctx context.Context, req *http.Request, ownerID uuid.UUID, statuses []entities.AppointmentStatus) paginate.Page
	FindAppointmentByAppCode(appCode string) (*entities.Appointment, error)
	FindAndLock(appCode string, tx *gorm.DB) (*entities.Appointment, error)
	FindByIDAndOwner(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*entities.Appointment, error)
	Update(appointment *entities.Appointment) error
	UpdateStatus(ctx context.Context, appointmentID uuid.UUID, status entities.AppointmentStatus) error
	MarkAppointmentsOngoing(ctx context.Context, now time.Time) (int64, error)
	MarkAppointmentsCompleted(ctx context.Context, now time.Time) (int64, error)
	MarkAppointmentsExpired(ctx context.Context, now time.Time) (int64, error)
	WithTx(tx *gorm.DB) AppointmentRepository
}

type gormAppointmentRepository struct {
	db *gorm.DB
}

func NewGormAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &gormAppointmentRepository{db: db}
}

func (r *gormAppointmentRepository) WithTx(tx *gorm.DB) AppointmentRepository {
	return &gormAppointmentRepository{db: tx}
}

func (r *gormAppointmentRepository) Create(appointment *entities.Appointment) error {
	if err := r.db.Create(appointment).Error; err != nil {
		return apperr.TranslateRepoError("repository.appointment.Create", err)
	}
	return nil
}

func (r *gormAppointmentRepository) GetAppointmentsByOwnerIDQuery(ctx context.Context, req *http.Request, ownerID uuid.UUID, statuses []entities.AppointmentStatus) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Appointment{}).Where("owner_id = ?", ownerID).Order("created_at DESC")
	if len(statuses) > 0 {
		values := make([]string, 0, len(statuses))
		for _, status := range statuses {
			values = append(values, string(status))
		}
		db = db.Where("status IN ?", values)
	}
	var request interface{}
	if req != nil {
		request = req
	} else {
		request = &paginate.Request{}
	}
	return pg.With(db).Request(request).Response(&[]responses.AppointmentResponse{})
}

func (r *gormAppointmentRepository) FindAppointmentByAppCode(appCode string) (*entities.Appointment, error) {
	var appointment entities.Appointment
	if err := r.db.Where("app_code = ?", appCode).First(&appointment).Error; err != nil {
		return nil, apperr.TranslateRepoError("repository.appointment.FindByAppCode", err)
	}
	return &appointment, nil
}

func (r *gormAppointmentRepository) FindAndLock(appCode string, tx *gorm.DB) (*entities.Appointment, error) {
	var appointment entities.Appointment
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("app_code = ?", appCode).First(&appointment).Error; err != nil {
		return nil, apperr.TranslateRepoError("repository.appointment.FindAndLock", err)
	}
	return &appointment, nil
}

func (r *gormAppointmentRepository) FindByIDAndOwner(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*entities.Appointment, error) {
	var appointment entities.Appointment
	if err := r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", id, ownerID).First(&appointment).Error; err != nil {
		return nil, apperr.TranslateRepoError("repository.appointment.FindByIDAndOwner", err)
	}
	return &appointment, nil
}

func (r *gormAppointmentRepository) Update(appointment *entities.Appointment) error {
	if err := r.db.Save(appointment).Error; err != nil {
		return apperr.TranslateRepoError("repository.appointment.Update", err)
	}
	return nil
}

func (r *gormAppointmentRepository) UpdateStatus(ctx context.Context, appointmentID uuid.UUID, status entities.AppointmentStatus) error {
	res := r.db.WithContext(ctx).Model(&entities.Appointment{}).
		Where("id = ?", appointmentID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return apperr.TranslateRepoError("repository.appointment.UpdateStatus", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperr.TranslateRepoError("repository.appointment.UpdateStatus", gorm.ErrRecordNotFound)
	}
	return nil
}

func (r *gormAppointmentRepository) MarkAppointmentsOngoing(ctx context.Context, now time.Time) (int64, error) {
	res := r.db.WithContext(ctx).Model(&entities.Appointment{}).
		Where("status = ?", entities.AppointmentStatusPending).
		Where("start_time <= ?", now).
		Updates(map[string]interface{}{
			"status":     entities.AppointmentStatusOngoing,
			"updated_at": now,
		})
	if res.Error != nil {
		return 0, apperr.TranslateRepoError("repository.appointment.MarkAppointmentsOngoing", res.Error)
	}
	return res.RowsAffected, nil
}

func (r *gormAppointmentRepository) MarkAppointmentsCompleted(ctx context.Context, now time.Time) (int64, error) {
	deadlineExpr := "date_trunc('day', end_date) + (end_time - date_trunc('day', end_time))"
	res := r.db.WithContext(ctx).Model(&entities.Appointment{}).
		Where("status IN ?", []entities.AppointmentStatus{entities.AppointmentStatusPending, entities.AppointmentStatusOngoing}).
		Where(deadlineExpr+" < ?", now).
		Where("attendees_booked > 0").
		Updates(map[string]interface{}{
			"status":     entities.AppointmentStatusCompleted,
			"updated_at": now,
		})
	if res.Error != nil {
		return 0, apperr.TranslateRepoError("repository.appointment.MarkAppointmentsCompleted", res.Error)
	}
	return res.RowsAffected, nil
}

func (r *gormAppointmentRepository) MarkAppointmentsExpired(ctx context.Context, now time.Time) (int64, error) {
	deadlineExpr := "date_trunc('day', end_date) + (end_time - date_trunc('day', end_time))"
	res := r.db.WithContext(ctx).Model(&entities.Appointment{}).
		Where("status IN ?", []entities.AppointmentStatus{entities.AppointmentStatusPending, entities.AppointmentStatusOngoing}).
		Where(deadlineExpr+" < ?", now).
		Where("attendees_booked = 0").
		Updates(map[string]interface{}{
			"status":     entities.AppointmentStatusExpired,
			"updated_at": now,
		})
	if res.Error != nil {
		return 0, apperr.TranslateRepoError("repository.appointment.MarkAppointmentsExpired", res.Error)
	}
	return res.RowsAffected, nil
}
