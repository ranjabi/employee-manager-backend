package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"

	"employee-manager/constants"
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
	pgConn := db.Setup(ctx)

	managerRepository := repositories.NewManagerRepository(ctx, pgConn)
	departmentRepository := repositories.NewDepartmentRepository(ctx, pgConn)
	employeeRepository := repositories.NewEmployeeRepository(ctx, pgConn)

	authService := services.NewAuthService(managerRepository)
	managerService := services.NewManagerService(managerRepository)
	departmentService := services.NewDepartmentService(departmentRepository)
	employeeService := services.NewEmployeeService(employeeRepository)

	authHandler := handlers.NewAuthHandler(authService)
	managerHandler := handlers.NewManagerHandler(managerService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/v1", func(r chi.Router) {
		// public
		r.Group(func(r chi.Router) {
			r.Post("/auth", AppHandler(authHandler.HandleRegisterLoginManager))
		})

		// protected
		r.Group(func(r chi.Router) {
			tokenAuth := jwtauth.New(constants.HASH_ALG, []byte(constants.JWT_SECRET), nil)
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.Get("/user", AppHandler(managerHandler.HandleGetProfile))
			r.Patch("/user", AppHandler(managerHandler.HandleUpdateProfile))

			r.Post("/department", AppHandler(departmentHandler.HandleCreateDepartment))
			r.Get("/department", AppHandler(departmentHandler.HandleGetAllDepartment))
			r.Patch("/department/{departmentId}", AppHandler(departmentHandler.HandleUpdateDepartment))
			r.Delete("/department/{departmentId}", AppHandler(departmentHandler.HandleDeleteDepartment))

			r.Post("/employee", AppHandler(employeeHandler.HandleCreateEmployee))
			r.Get("/employee", AppHandler(employeeHandler.HandleGetAllEmployee))
			r.Patch("/employee/{identityNumber}", AppHandler(employeeHandler.HandleUpdateEmployee))
		})
	})

	http.ListenAndServe(":8080", r)
}
