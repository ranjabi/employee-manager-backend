package services

import (
	"employee-manager/models"
	"employee-manager/repositories"
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