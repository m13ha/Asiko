package repository

import (
    "github.com/google/uuid"
    apperr "github.com/m13ha/asiko/errors"
    "github.com/m13ha/asiko/models/entities"
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
    if err := r.db.Create(entry).Error; err != nil {
        return apperr.TranslateRepoError("repository.banlist.Create", err)
    }
    return nil
}

func (r *gormBanListRepository) Delete(userID uuid.UUID, email string) error {
    if err := r.db.Where("user_id = ? AND banned_email = ?", userID, email).Delete(&entities.BanListEntry{}).Error; err != nil {
        return apperr.TranslateRepoError("repository.banlist.Delete", err)
    }
    return nil
}

func (r *gormBanListRepository) FindByUserAndEmail(userID uuid.UUID, email string) (*entities.BanListEntry, error) {
	var entry entities.BanListEntry
    if err := r.db.Where("user_id = ? AND banned_email = ?", userID, email).First(&entry).Error; err != nil {
        return nil, apperr.TranslateRepoError("repository.banlist.FindByUserAndEmail", err)
    }
    return &entry, nil
}

func (r *gormBanListRepository) GetAllByUser(userID uuid.UUID) ([]entities.BanListEntry, error) {
	var entries []entities.BanListEntry
    if err := r.db.Where("user_id = ?", userID).Find(&entries).Error; err != nil {
        return nil, apperr.TranslateRepoError("repository.banlist.GetAllByUser", err)
    }
    return entries, nil
}
