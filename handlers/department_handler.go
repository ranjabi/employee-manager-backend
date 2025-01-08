package handlers

import (
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type DepartmentHandler struct {
	departmentService services.DepartmentService
}

func NewDepartmentHandler(departmentService services.DepartmentService) DepartmentHandler {
	return DepartmentHandler{departmentService}
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
