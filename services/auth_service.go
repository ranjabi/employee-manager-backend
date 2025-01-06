package services

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"employee-manager/models"
	"employee-manager/repositories"
)

type AuthService struct {
	managerRepository repositories.ManagerRepository
}

func NewAuthService(managerRepository repositories.ManagerRepository) AuthService {
	return AuthService{managerRepository}
}

func (s *AuthService) CreateManager(manager models.Manager) (*models.Manager, error) {
	newManager, err := s.managerRepository.CreateManager(manager)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, models.NewError(http.StatusConflict, "Email is already taken")
		}
		return nil, err
	}

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	claims := map[string]any{
		"manager_id":    newManager.Id,
		"manager_email": newManager.Email,
	}
	_, tokenString, _ := tokenAuth.Encode(claims)

	newManager.Token = tokenString

	return newManager, nil
}
