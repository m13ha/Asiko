package services

import (
	"github.com/google/uuid"
	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/repository"
	"github.com/m13ha/appointment_master/utils"
)

type BanListService interface {
	AddToBanList(userID uuid.UUID, email string) (*entities.BanListEntry, error)
	RemoveFromBanList(userID uuid.UUID, email string) error
	GetBanList(userID uuid.UUID) ([]entities.BanListEntry, error)
}

type banListServiceImpl struct {
	banListRepo repository.BanListRepository
}

func NewBanListService(banListRepo repository.BanListRepository) BanListService {
	return &banListServiceImpl{banListRepo: banListRepo}
}

func (s *banListServiceImpl) AddToBanList(userID uuid.UUID, email string) (*entities.BanListEntry, error) {
	normalizedEmail := utils.NormalizeEmail(email)
	_, err := s.banListRepo.FindByUserAndEmail(userID, normalizedEmail)
	if err == nil {
		return nil, myerrors.NewUserError("email already on ban list")
	}

	entry := &entities.BanListEntry{
		UserID:      userID,
		BannedEmail: normalizedEmail,
	}

	if err := s.banListRepo.Create(entry); err != nil {
		return nil, myerrors.NewUserError("failed to add email to ban list")
	}

	return entry, nil
}

func (s *banListServiceImpl) RemoveFromBanList(userID uuid.UUID, email string) error {
	normalizedEmail := utils.NormalizeEmail(email)
	return s.banListRepo.Delete(userID, normalizedEmail)
}

func (s *banListServiceImpl) GetBanList(userID uuid.UUID) ([]entities.BanListEntry, error) {
	return s.banListRepo.GetAllByUser(userID)
}
