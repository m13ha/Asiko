package repository

import (
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"gorm.io/gorm"
)

type BanListRepository interface {
	Create(entry *entities.BanListEntry) error
	Delete(userID uuid.UUID, email string) error
	FindByUserAndEmail(userID uuid.UUID, email string) (*entities.BanListEntry, error)
	GetAllByUser(userID uuid.UUID) ([]entities.BanListEntry, error)
}

type gormBanListRepository struct {
	db *gorm.DB
}

func NewGormBanListRepository(db *gorm.DB) BanListRepository {
	return &gormBanListRepository{db: db}
}

func (r *gormBanListRepository) Create(entry *entities.BanListEntry) error {
	return r.db.Create(entry).Error
}

func (r *gormBanListRepository) Delete(userID uuid.UUID, email string) error {
	return r.db.Where("user_id = ? AND banned_email = ?", userID, email).Delete(&entities.BanListEntry{}).Error
}

func (r *gormBanListRepository) FindByUserAndEmail(userID uuid.UUID, email string) (*entities.BanListEntry, error) {
	var entry entities.BanListEntry
	if err := r.db.Where("user_id = ? AND banned_email = ?", userID, email).First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *gormBanListRepository) GetAllByUser(userID uuid.UUID) ([]entities.BanListEntry, error) {
	var entries []entities.BanListEntry
	if err := r.db.Where("user_id = ?", userID).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}
