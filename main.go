package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var s3Client *s3.Client

func initS3(ctx context.Context) error {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(awsConfig)
	return nil
}

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
	if err := initS3(ctx); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Running migration...")
	migrate, err := migrate.New(
		"file://db/migrations",
		db.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error connecting to db:" + err.Error())
	}
	if err := migrate.Up(); err != nil {
		if err.Error() == "no change" {
			fmt.Println("no change")
		} else {
			log.Fatal("Migration failed:" + err.Error())
		}
	}
	fmt.Println("Migration success")

	managerRepository := repositories.NewManagerRepository(ctx, pgConn)
	departmentRepository := repositories.NewDepartmentRepository(ctx, pgConn)
	employeeRepository := repositories.NewEmployeeRepository(ctx, pgConn)

	authService := services.NewAuthService(managerRepository)
	managerService := services.NewManagerService(managerRepository)
	departmentService := services.NewDepartmentService(departmentRepository)
	employeeService := services.NewEmployeeService(employeeRepository)
	fileService := services.NewFileService(s3Client, ctx)

	authHandler := handlers.NewAuthHandler(authService)
	managerHandler := handlers.NewManagerHandler(managerService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	fileHandler := handlers.NewFileHandler(fileService)

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
			r.Use(AllowContentType("application/json", "multipart/form-data"))

			r.Get("/user", AppHandler(managerHandler.HandleGetProfile))
			r.Patch("/user", AppHandler(managerHandler.HandleUpdateProfile))

			r.Post("/department", AppHandler(departmentHandler.HandleCreateDepartment))
			r.Get("/department", AppHandler(departmentHandler.HandleGetAllDepartment))
			r.Patch("/department/{departmentId}", AppHandler(departmentHandler.HandleUpdateDepartment))
			r.Delete("/department/{departmentId}", AppHandler(departmentHandler.HandleDeleteDepartment))

			r.Post("/employee", AppHandler(employeeHandler.HandleCreateEmployee))
			r.Get("/employee", AppHandler(employeeHandler.HandleGetAllEmployee))
			r.Patch("/employee/{identityNumber}", AppHandler(employeeHandler.HandleUpdateEmployee))
			r.Delete("/employee/{identityNumber}", AppHandler(employeeHandler.HandleDeleteEmployee))

			r.Post("/file", AppHandler(fileHandler.HandleUploadFile))
		})
	})

	http.ListenAndServe(":8080", r)
}

func AllowContentType(contentTypes ...string) func(http.Handler) http.Handler {
	allowedContentTypes := make(map[string]struct{}, len(contentTypes))
	for _, ctype := range contentTypes {
		allowedContentTypes[strings.TrimSpace(strings.ToLower(ctype))] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				// Skip check for empty content body
				next.ServeHTTP(w, r)
				return
			}

			s := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))

			if _, ok := allowedContentTypes[s]; ok {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
		})
	}
}