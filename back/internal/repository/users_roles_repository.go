package repository

import (
	"backend-app/internal/domain"
	"gorm.io/gorm"
)

type UsersRolesRepository struct {
	db *gorm.DB
}

func NewUsersRolesRepository(db *gorm.DB) *UsersRolesRepository {
	return &UsersRolesRepository{db: db}
}

func (r *UsersRolesRepository) AssignRoleToUser(user *domain.User, role *domain.Role) error {
	return r.db.Model(user).Association("Roles").Append(role)
}

func (r *UsersRolesRepository) RemoveRoleFromUser(user *domain.User, role *domain.Role) error {
	return r.db.Model(user).Association("Roles").Delete(role)
}

func (r *UsersRolesRepository) GetUserRoles(user *domain.User) ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.db.Model(user).Association("Roles").Find(&roles)
	return roles, err
}