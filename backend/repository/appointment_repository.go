package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppointmentRepository interface {
	Create(appointment *entities.Appointment) error
	GetAppointmentsByOwnerIDQuery(ctx context.Context, ownerID uuid.UUID) paginate.Page
	FindAppointmentByAppCode(appCode string) (*entities.Appointment, error)
	FindAndLock(appCode string, tx *gorm.DB) (*entities.Appointment, error)
	Update(appointment *entities.Appointment) error
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
	return r.db.Create(appointment).Error
}

func (r *gormAppointmentRepository) GetAppointmentsByOwnerIDQuery(ctx context.Context, ownerID uuid.UUID) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Appointment{}).Where("owner_id = ?", ownerID).Order("created_at DESC")
	return pg.With(db).Request(ctx).Response(&[]responses.AppointmentResponse{})
}

func (r *gormAppointmentRepository) FindAppointmentByAppCode(appCode string) (*entities.Appointment, error) {
	var appointment entities.Appointment
	if err := r.db.Where("app_code = ?", appCode).First(&appointment).Error; err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *gormAppointmentRepository) FindAndLock(appCode string, tx *gorm.DB) (*entities.Appointment, error) {
	var appointment entities.Appointment
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("app_code = ?", appCode).First(&appointment).Error; err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *gormAppointmentRepository) Update(appointment *entities.Appointment) error {
	return r.db.Save(appointment).Error
}
