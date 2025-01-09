package handlers

import (
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type EmployeeHandler struct {
	employeeService services.EmployeeService
}

func NewEmployeeHandler(employeeService services.EmployeeService) EmployeeHandler {
	return EmployeeHandler{employeeService}
}

func (h *EmployeeHandler) HandleCreateEmployee(w http.ResponseWriter, r *http.Request) error {
	payload := struct {
		IdentityNumber   string `json:"identityNumber" validate:"required,min=5,max=33"`
		Name             string `json:"name" validate:"required,min=4,max=33"`
		EmployeeImageUri string `json:"employeeImageUri" validate:"required,uri"`
		Gender           string `json:"gender" validate:"required,oneof=male female"`
		DepartmentId     string `json:"departmentId" validate:"required"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return models.NewError(http.StatusBadRequest, err.Error())
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErr := fmt.Errorf("validation for '%s' failed", err.Field())
			return models.NewError(http.StatusBadRequest, validationErr.Error())
		}
	}

	newEmployee, err := h.employeeService.CreateEmployee(models.Employee{
		IdentityNumber:   payload.IdentityNumber,
		Name:             payload.Name,
		EmployeeImageUri: payload.EmployeeImageUri,
		Gender:           payload.Gender,
		DepartmentId:     payload.DepartmentId,
	})
	if err != nil {
		return err
	}

	lib.SetJsonResponse(w, http.StatusCreated)
	err = json.NewEncoder(w).Encode(newEmployee)
	if err != nil {
		return err
	}

	return nil
}

func (h *EmployeeHandler) HandleGetAllEmployee(w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()	
	identityNumber := params.Get("identityNumber")
	name := params.Get("name")
	gender := params.Get("gender")
	departmentId := params.Get("departmentId")
	limitStr := params.Get("limit")
	offsetStr := params.Get("offset")
	limit := 5
	offset := 0
	if limitStr != "" {
		limitTemp, err := strconv.Atoi(limitStr); 
		if err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
		if limitTemp >= 0 {
			limit = limitTemp
		}
	}
	if offsetStr != "" {
		offsetTemp, err := strconv.Atoi(offsetStr)
		if err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
		if offsetTemp >= 0 {
			offset = offsetTemp
		}
	}

	employees, err := h.employeeService.GetAllEmployee(offset, limit, identityNumber, name, gender, departmentId)
	if err != nil {
		return err
	}

	lib.SetJsonResponse(w, http.StatusOK)
	err = json.NewEncoder(w).Encode(employees)
	if err != nil {
		return err
	}
	
	return nil
}
