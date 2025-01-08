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

type DepartmentHandler struct {
	departmentService services.DepartmentService
}

func NewDepartmentHandler(departmentService services.DepartmentService) DepartmentHandler {
	return DepartmentHandler{departmentService}
}

func (h *DepartmentHandler) HandleGetAllDepartment(w http.ResponseWriter, r *http.Request) error {
	var err error
	limit := 5
	offset := 0
	params := r.URL.Query()	
	limitStr := params.Get("limit")
	offsetStr := params.Get("offset")
	name := params.Get("name")
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
	
	departments, err := h.departmentService.GetAllDepartment(offset, limit, name)
	if err != nil {
		return err
	}
	lib.SetJsonResponse(w, http.StatusOK)
	err = json.NewEncoder(w).Encode(departments)
	if err != nil {
		return err
	}
	
	return nil
}

func (h *DepartmentHandler) HandleUpdateDepartment(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	payload := types.UpdateDepartmentProfilePayload{}
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

	department, err := h.departmentService.PartialUpdate(id, payload)
	if err != nil {
		return err
	}

	res := struct {
		DepartmentId string `json:"departmentId"`
		Name         string `json:"name"`
	}{
		DepartmentId: department.Id,
		Name: department.Name,
	}
	lib.SetJsonResponse(w, http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}

	return nil
}

func (h *DepartmentHandler) HandleCreateDepartment(w http.ResponseWriter, r *http.Request) error {
	payload := struct {
		Name string `json:"name" validate:"required,min=4,max=33"`
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

	newDepartment, err := h.departmentService.CreateDepartment(models.Department{
		Name: payload.Name,
	})
	if err != nil {
		return err
	}

	res := struct {
		DepartmentId string `json:"departmentId"`
		Name         string `json:"name"`
	}{
		DepartmentId: newDepartment.Id,
		Name: newDepartment.Name,
	}
	lib.SetJsonResponse(w, http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}

	return nil
}