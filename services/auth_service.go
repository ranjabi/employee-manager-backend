package services

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"employee-manager/models"
	"employee-manager/repositories"
)

var saltRound int = 10
var jwtSecret string = "secret"
var hashAlg string = "HS256"
var uniqueViolationErrorCode string = "23505"

type AuthService struct {
	managerRepository repositories.ManagerRepository
}

func NewAuthService(managerRepository repositories.ManagerRepository) AuthService {
	return AuthService{managerRepository}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), saltRound)
	return string(bytes), err
}

func CreateClaims(manager *models.Manager) (string, error) {
	tokenAuth := jwtauth.New(hashAlg, []byte(jwtSecret), nil)
	claims := map[string]any{
		"manager_id":    manager.Id,
		"manager_email": manager.Email,
	}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) CreateManager(manager models.Manager) (*models.Manager, error) {
	hashedPassword, err := HashPassword(manager.Password)
	if err != nil {
		return nil, err
	}
	manager.Password = hashedPassword

	newManager, err := s.managerRepository.Save(manager)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == uniqueViolationErrorCode {
			return nil, models.NewError(http.StatusConflict, "Email is already taken")
		}
		return nil, err
	}

	token, err := CreateClaims(newManager)
	if err != nil {
		return nil, err
	}

	newManager.Token = token

	return newManager, nil
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}

func (s *AuthService) Login(email string, password string) (*models.Manager, error) {
	manager, err := s.managerRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewError(http.StatusNotFound, "Email is not exist")
		}

		return nil, err
	}

	match := CheckPasswordHash(manager.Password, password)
	if match {
		token, err := CreateClaims(manager)
		if err != nil {
			return nil, err
		}

		manager.Token = token
		return manager, nil
	} else {
		return nil, models.NewError(http.StatusUnauthorized, "Invalid email/password")
	}
}
