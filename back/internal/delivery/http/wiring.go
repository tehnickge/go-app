package httpdelivery

import (
	"backend-app/internal/usecase"
)

type Dependencies struct {
	Auth        *usecase.AuthUsecase
	Users       *usecase.UsersUsecase
	Roles       *usecase.RolesUsecase
	Groups      *usecase.GroupsUsecase
	Relations   *usecase.RelationsUsecase

	JWTSecret string
}

