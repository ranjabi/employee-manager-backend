package handlers

import (
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/services"
	"employee-manager/types"
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
	if !lib.IsValidURI(payload.EmployeeImageUri) {
		return models.NewError(http.StatusBadRequest, "invalid uri")
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

func (h *EmployeeHandler) HandleUpdateEmployee(w http.ResponseWriter, r *http.Request) error {
	// TODO none request body result in EOF, at least there must be { }
	identityNumber := r.PathValue("identityNumber")
	payload := types.UpdateEmployeePayload{}
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

	employee, err := h.employeeService.PartialUpdate(identityNumber, payload)
	if err != nil {
		return err
	}

	res := struct {
		IdentityNumber string `json:"identityNumber"`
		Name string `json:"name"`
		EmployeeImageUri string `json:"employeeImageUri"`
		Gender string `json:"gender"`
		DepartmentId string `json:"departmentId"`
	}{
		IdentityNumber: employee.IdentityNumber,
		Name: employee.Name,
		EmployeeImageUri: employee.EmployeeImageUri,
		Gender: employee.Gender,
		DepartmentId: employee.DepartmentId,
	}
	lib.SetJsonResponse(w, http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	
	return nil
}

func (h *EmployeeHandler) HandleDeleteEmployee(w http.ResponseWriter, r *http.Request) error {
	identityNumber := r.PathValue("identityNumber")

	err := h.employeeService.Delete(identityNumber)
	if err != nil {
		return err
	}

	w.Write([]byte(""))

	return nil
}