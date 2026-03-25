package repository

import (
	"backend-app/internal/domain"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *RoleRepository) List() ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.db.Find(&roles).Error
	return roles, err
}