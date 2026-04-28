package main

import (
	"log"
	"net/http"
	"time"

	"backend-app/internal/config"
	"backend-app/internal/migrations"
	"backend-app/internal/repository"
	httpdelivery "backend-app/internal/delivery/http"
	"backend-app/internal/usecase"
	"backend-app/pkg/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Загружаем конфиг
	cfg := config.Load()

	// Подключаемся к базе
	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем миграции
	migrations.RunMigrations(database)

	// Repositories
	userRepo := repository.NewUserRepository(database)
	roleRepo := repository.NewRoleRepository(database)
	groupRepo := repository.NewGroupRepository(database)
	urRepo := repository.NewUsersRolesRepository(database)
	ugRepo := repository.NewUsersGroupsRepository(database)

	jwtTTL, err := time.ParseDuration(cfg.JWT.TTL)
	if err != nil {
		log.Fatal("invalid JWT_TTL:", err)
	}

	// Usecases
	authz := usecase.NewAuthorizer(urRepo)
	authUC := usecase.NewAuthUsecase(userRepo, roleRepo, urRepo, cfg.JWT.Secret, jwtTTL)
	usersUC := usecase.NewUsersUsecase(userRepo, authz)
	rolesUC := usecase.NewRolesUsecase(roleRepo, authz)
	groupsUC := usecase.NewGroupsUsecase(groupRepo, authz)
	relationsUC := usecase.NewRelationsUsecase(urRepo, ugRepo, authz)

	deps := httpdelivery.Dependencies{
		Auth:        authUC,
		Users:       usersUC,
		Roles:       rolesUC,
		Groups:      groupsUC,
		Relations:   relationsUC,
		JWTSecret:   cfg.JWT.Secret,
	}

	// Router
	r := chi.NewRouter()
	httpdelivery.RegisterRoutes(r, deps)

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", r)
}