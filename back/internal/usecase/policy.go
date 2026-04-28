package usecase

import (
	"backend-app/internal/repository"

	"github.com/google/uuid"
)

type Authorizer struct {
	urRepo *repository.UsersRolesRepository
}

func NewAuthorizer(urRepo *repository.UsersRolesRepository) *Authorizer {
	return &Authorizer{urRepo: urRepo}
}

func (a *Authorizer) IsAdmin(callerID uuid.UUID) (bool, error) {
	roles, err := a.urRepo.GetActiveRolesByUserID(callerID)
	if err != nil {
		return false, err
	}
	for _, r := range roles {
		if r != nil && r.Name == "admin" {
			return true, nil
		}
	}
	return false, nil
}

