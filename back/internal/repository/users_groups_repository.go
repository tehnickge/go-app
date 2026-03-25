package repository

import (
	"backend-app/internal/domain"
	"gorm.io/gorm"
)

type UsersGroupsRepository struct {
	db *gorm.DB
}

func NewUsersGroupsRepository(db *gorm.DB) *UsersGroupsRepository {
	return &UsersGroupsRepository{db: db}
}

func (r *UsersGroupsRepository) AddUserToGroup(user *domain.User, group *domain.Group, role string) error {
	return r.db.Model(group).Association("Users").Append(user)
}

func (r *UsersGroupsRepository) RemoveUserFromGroup(user *domain.User, group *domain.Group) error {
	return r.db.Model(group).Association("Users").Delete(user)
}