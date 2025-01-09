package services

import (
	"employee-manager/constants"
	"employee-manager/models"
	"employee-manager/repositories"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type EmployeeService struct {
	employeeRepository repositories.EmployeeRepository
}

func NewEmployeeService(employeeRepository repositories.EmployeeRepository) EmployeeService {
	return EmployeeService{employeeRepository}
}

func (s *EmployeeService) CreateEmployee(employee models.Employee) (*models.Employee, error) {
	newEmployee, err := s.employeeRepository.Save(employee)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.FOREIGN_KEY_CONSTRAINT_VIOLATION_ERROR_CODE {
			return nil, models.NewError(http.StatusBadRequest, "Invalid department id")
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.UNIQUE_VIOLATION_ERROR_CODE {
			return nil, models.NewError(http.StatusConflict, "Identity number is already taken")
		}
		return nil, err
	}

	return newEmployee, nil
}