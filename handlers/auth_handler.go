package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return AuthHandler{authService}
}

func (h *AuthHandler) HandleRegisterLoginManager(w http.ResponseWriter, r *http.Request) error {
	payload := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"`
		Action   string `json:"action" validate:"required,oneof=create login"`
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

	switch payload.Action {
	case "create":
		newManager, err := h.authService.CreateManager(models.Manager{
			Email:    payload.Email,
			Password: payload.Password,
		})
		if err != nil {
			return err
		}

		response := struct {
			Email string `json:"email"`
			Token string `json:"token"`
		}{
			Email: newManager.Email,
			Token: newManager.Token,
		}
		lib.SetJsonResponse(w, http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return err
		}
	case "login":
		manager, err := h.authService.Login(payload.Email, payload.Password)
		if err != nil {
			return err
		}

		response := struct {
			Email string `json:"email"`
			Token string `json:"token"`
		}{
			Email: manager.Email,
			Token: manager.Token,
		}
		lib.SetJsonResponse(w, http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return err
		}
	}

	return nil
}
