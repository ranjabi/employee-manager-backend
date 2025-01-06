package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"employee-manager/db"
	"employee-manager/handlers"
	"employee-manager/models"
	"employee-manager/repositories"
	"employee-manager/services"
)

func AppHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if err, ok := err.(*models.AppError); ok {
				if err.Code != 0 {
					http.Error(w, err.Error(), err.Code)
					return
				}
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	r := chi.NewRouter()
	pgConn := db.Setup(ctx)

	managerRepository := repositories.NewManagerRepository(ctx, pgConn)
	authService := services.NewAuthService(managerRepository)
	authHandler := handlers.NewAuthHandler(authService)

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/auth", AppHandler(authHandler.HandleRegisterLoginManager))
	})

	http.ListenAndServe(":8080", r)
}
