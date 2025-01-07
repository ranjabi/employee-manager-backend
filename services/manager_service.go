package services

import (
	"employee-manager/constants"
	"employee-manager/models"
	"employee-manager/repositories"
	"employee-manager/types"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ManagerService struct {
	managerRepository repositories.ManagerRepository
}

func NewManagerService(managerRepository repositories.ManagerRepository) ManagerService {
	return ManagerService{managerRepository}
}

func (s *ManagerService) FindById(id string) (*models.Manager, error) {
	manager, err := s.managerRepository.FindById(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewError(http.StatusNotFound, "")
		}

		return nil, err
	}

	return manager, nil
}

func (s *ManagerService) PartialUpdate(id string, payload types.UpdateManagerProfilePayload) (*models.Manager, error) {
	manager, err := s.managerRepository.PartialUpdate(id, payload)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.UNIQUE_VIOLATION_ERROR_CODE {
			return nil, models.NewError(http.StatusConflict, "Email is already taken")
		}
		return nil, err
	}

	return manager, nil
}