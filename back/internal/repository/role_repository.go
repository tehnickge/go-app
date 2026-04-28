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

func (r *RoleRepository) GetByID(id string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&role).Error
	return &role, err
}

func (r *RoleRepository) UpdateFields(roleID string, fields map[string]interface{}) error {
	return r.db.Model(&domain.Role{}).Where("id = ? AND deleted_at IS NULL", roleID).Updates(fields).Error
}

func (r *RoleRepository) SoftDelete(roleID string) error {
	// Uses direct deleted_at update to keep API of repositories simple.
	return r.db.Model(&domain.Role{}).Where("id = ? AND deleted_at IS NULL", roleID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}