package services

import (
    "github.com/google/uuid"
    serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
    "github.com/m13ha/asiko/models/entities"
    "github.com/m13ha/asiko/repository"
    "github.com/m13ha/asiko/utils"
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
        return nil, serviceerrors.ConflictError("email already on ban list")
    }

	entry := &entities.BanListEntry{
		UserID:      userID,
		BannedEmail: normalizedEmail,
	}

    if err := s.banListRepo.Create(entry); err != nil {
        return nil, serviceerrors.FromError(err)
    }

	return entry, nil
}

func (s *banListServiceImpl) RemoveFromBanList(userID uuid.UUID, email string) error {
    normalizedEmail := utils.NormalizeEmail(email)
    if err := s.banListRepo.Delete(userID, normalizedEmail); err != nil {
        return serviceerrors.FromError(err)
    }
    return nil
}

func (s *banListServiceImpl) GetBanList(userID uuid.UUID) ([]entities.BanListEntry, error) {
    entries, err := s.banListRepo.GetAllByUser(userID)
    if err != nil {
        return nil, serviceerrors.FromError(err)
    }
    return entries, nil
}
