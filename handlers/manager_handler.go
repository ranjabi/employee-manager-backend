package handlers

import (
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/services"
	"employee-manager/types"
	"encoding/json"
	"fmt"
	"io"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.NewError(http.StatusBadRequest, err.Error())
	}

	payload := types.UpdateManagerProfilePayload{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return models.NewError(http.StatusBadRequest, err.Error())
	}

	// TODO rewrite using for loop
	if string(payload.EmailRaw) == "null" ||
		string(payload.NameRaw) == "null" ||
		string(payload.UserImageUriRaw) == "null" ||
		string(payload.CompanyNameRaw) == "null" ||
		string(payload.CompanyImageUriRaw) == "null" {
		return models.NewError(http.StatusBadRequest, "input can't be null")
	}

	if payload.EmailRaw != nil {
		if err := json.Unmarshal([]byte(payload.EmailRaw), &payload.Email); err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
	}
	if payload.NameRaw != nil {
		if err := json.Unmarshal([]byte(payload.NameRaw), &payload.Name); err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
	}
	if payload.UserImageUriRaw != nil {
		if err := json.Unmarshal([]byte(payload.UserImageUriRaw), &payload.UserImageUri); err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
		if !lib.IsValidURI(*payload.UserImageUri) {
			return models.NewError(http.StatusBadRequest, "invalid uri")
		}
	}
	if payload.CompanyNameRaw != nil {
		if err := json.Unmarshal([]byte(payload.CompanyNameRaw), &payload.CompanyName); err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
	}
	if payload.CompanyImageUriRaw != nil {
		if err := json.Unmarshal([]byte(payload.CompanyImageUriRaw), &payload.CompanyImageUri); err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
		if !lib.IsValidURI(*payload.CompanyImageUri) {

			return models.NewError(http.StatusBadRequest, "invalid uri")
		}
	}

	if err := validator.New().Struct(payload); err != nil {
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
