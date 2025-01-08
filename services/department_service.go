package services

import (
	"employee-manager/models"
	"employee-manager/repositories"
	"employee-manager/types"
)

type DepartmentService struct {
	departmentRepository repositories.DepartmentRepository
}

func NewDepartmentService(departmentRepository repositories.DepartmentRepository) DepartmentService {
	return DepartmentService{departmentRepository}
}

func (s *DepartmentService) GetAllDepartment(offset int, limit int, name string) ([]models.Department, error) {
	departments, err := s.departmentRepository.GetAllDepartment(offset, limit, name)
	if err != nil {
		return nil, err
	}

	return departments, nil
}

func (s *DepartmentService) CreateDepartment(department models.Department) (*models.Department, error) {
	newDepartment, err := s.departmentRepository.Save(department)
	if err != nil {
		return nil, err
	}

	return newDepartment, nil
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
		return err
	}

	return nil
}
