package services

import (
	"employee-manager/constants"
	"employee-manager/models"
	"employee-manager/repositories"
	"employee-manager/types"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type DepartmentService struct {
	departmentRepository repositories.DepartmentRepository
}

func NewDepartmentService(departmentRepository repositories.DepartmentRepository) DepartmentService {
	return DepartmentService{departmentRepository}
}

func (s *DepartmentService) CreateDepartment(department models.Department) (*models.Department, error) {
	newDepartment, err := s.departmentRepository.Save(department)
	if err != nil {
		return nil, err
	}

	return newDepartment, nil
}

func (s *DepartmentService) GetAllDepartment(offset int, limit int, name string, managerId string) ([]models.Department, error) {
	departments, err := s.departmentRepository.GetAllDepartment(offset, limit, name, managerId)
	if err != nil {
		return nil, err
	}

	return departments, nil
}

func (s *DepartmentService) PartialUpdate(id string, payload types.UpdateDepartmentProfilePayload) (*models.Department, error) {
	department, err := s.departmentRepository.PartialUpdate(id, payload)
	if err != nil {
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) Delete(id string) error {
	err := s.departmentRepository.Delete(id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.FOREIGN_KEY_CONSTRAINT_VIOLATION_ERROR_CODE {
			return models.NewError(http.StatusConflict, "Still contain employee")
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.INVALID_INPUT_SYNTAX_TYPE_ERROR_CODE {
			return models.NewError(http.StatusNotFound, "")
		}
		
		return err
	}

	return nil
}
