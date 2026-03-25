package repository

import (
	"backend-app/internal/domain"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *domain.Group) error {
	return r.db.Create(group).Error
}

func (r *GroupRepository) GetByID(id string) (*domain.Group, error) {
	var group domain.Group
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&group).Error
	return &group, err
}

func (r *GroupRepository) SoftDelete(id string) error {
	return r.db.Model(&domain.Group{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}