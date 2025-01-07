package services

import (
	"employee-manager/models"
	"employee-manager/repositories"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
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
