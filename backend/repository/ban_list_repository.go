package repository

import (
    "github.com/google/uuid"
    repoerrors "github.com/m13ha/asiko/errors/repoerrors"
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
        return repoerrors.InternalError("failed to create ban list entry: " + err.Error())
    }
    return nil
}

func (r *gormBanListRepository) Delete(userID uuid.UUID, email string) error {
    if err := r.db.Where("user_id = ? AND banned_email = ?", userID, email).Delete(&entities.BanListEntry{}).Error; err != nil {
        return repoerrors.InternalError("failed to delete ban list entry: " + err.Error())
    }
    return nil
}

func (r *gormBanListRepository) FindByUserAndEmail(userID uuid.UUID, email string) (*entities.BanListEntry, error) {
	var entry entities.BanListEntry
    if err := r.db.Where("user_id = ? AND banned_email = ?", userID, email).First(&entry).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, repoerrors.NotFoundError("entry not found for user and email")
        }
        return nil, repoerrors.InternalError("failed to find ban list entry: " + err.Error())
    }
    return &entry, nil
}

func (r *gormBanListRepository) GetAllByUser(userID uuid.UUID) ([]entities.BanListEntry, error) {
	var entries []entities.BanListEntry
    if err := r.db.Where("user_id = ?", userID).Find(&entries).Error; err != nil {
        return nil, repoerrors.InternalError("failed to get all ban list entries for user: " + err.Error())
    }
    return entries, nil
}
