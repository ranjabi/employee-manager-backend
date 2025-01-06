package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"employee-manager/models"
	"employee-manager/services"
)

type Handler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) Handler {
	return Handler{authService}
}

func (h *Handler) HandleRegisterLoginManager(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"`
		Action   string `json:"action" validate:"required,oneof=create login"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErr := fmt.Errorf("validation for '%s' failed", err.Field())
			return models.NewError(http.StatusBadRequest, validationErr.Error())
		}
	}

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

	res, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}
