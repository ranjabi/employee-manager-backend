package handlers

import (
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
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
