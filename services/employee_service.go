package services

import (
	"employee-manager/constants"
	"employee-manager/models"
	"employee-manager/repositories"
	"employee-manager/types"
	"net/http"

	"github.com/google/uuid"
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

func (s *EmployeeService) GetAllEmployee(offset int, limit int, identityNumber string, name string, gender string, departmentId string) ([]models.Employee, error) {
	if gender != "male" && gender != "female" {
		return []models.Employee{}, nil
	}
	if err := uuid.Validate(departmentId); err != nil {
		return []models.Employee{}, nil
	}

	employees, err := s.employeeRepository.GetAllEmployee(offset, limit, identityNumber, name, gender, departmentId)
	if err != nil {
		return nil, err
	}

	return employees, err
}

func (s *EmployeeService) PartialUpdate(identityNumber string, payload types.UpdateEmployeePayload) (*models.Employee, error) {
	employee, err := s.employeeRepository.PartialUpdate(identityNumber, payload)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.FOREIGN_KEY_CONSTRAINT_VIOLATION_ERROR_CODE {
			return nil, models.NewError(http.StatusBadRequest, "Invalid department id")
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.INVALID_INPUT_SYNTAX_TYPE_ERROR_CODE {
			return nil, models.NewError(http.StatusBadRequest, "Invalid department id format")
		}

		return nil, err
	}

	return employee, nil
}

func (s *EmployeeService) Delete(identityNumber string) error {
	err := s.employeeRepository.Delete(identityNumber)
	if err != nil {
		return err
	}

	return nil
}