package handlers

import (
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/services"
	"employee-manager/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type ManagerHandler struct {
	managerService services.ManagerService
}

func NewManagerHandler(managerService services.ManagerService) ManagerHandler {
	return ManagerHandler{managerService}
}

func (h *ManagerHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) error {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	manager, err := h.managerService.FindById(claims["manager_id"].(string))
	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	response := struct {
		Email           string `json:"email"`
		Name            string `json:"name"`
		UserImageUri    string `json:"userImageUri"`
		CompanyName     string `json:"companyName"`
		CompanyImageUri string `json:"companyImageUri"`
	}{
		Email:           manager.Email,
		Name:            manager.Name.String,
		UserImageUri:    manager.UserImageUri.String,
		CompanyName:     manager.CompanyName.String,
		CompanyImageUri: manager.CompanyImageUri.String,
	}
	lib.SetJsonResponse(w, http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}

	return nil
}

func (h *ManagerHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) error {
	payload := types.UpdateManagerProfilePayload{}
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

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}
	managerId := claims["manager_id"].(string)
	manager, err := h.managerService.PartialUpdate(managerId, payload)
	if err != nil {
		return err
	}

	response := struct {
		Email           string `json:"email"`
		Name            string `json:"name"`
		UserImageUri    string `json:"userImageUri"`
		CompanyName     string `json:"companyName"`
		CompanyImageUri string `json:"companyImageUri"`
	}{
		Email:           manager.Email,
		Name:            manager.Name.String,
		UserImageUri:    manager.UserImageUri.String,
		CompanyName:     manager.CompanyName.String,
		CompanyImageUri: manager.CompanyImageUri.String,
	}
	lib.SetJsonResponse(w, http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}
	
	return nil
}


