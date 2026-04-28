package repository

import (
	"backend-app/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ? AND status != ?", email, "deleted").First(&user).Error
	return &user, err
}

func (r *UserRepository) SoftDelete(userID string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"status":     "deleted",
			"deleted_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ? AND status != ?", id, "deleted").First(&user).Error
	return &user, err
}

func (r *UserRepository) List() ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Where("status != ?", "deleted").Find(&users).Error
	return users, err
}

func (r *UserRepository) UpdateFields(userID string, fields map[string]interface{}) error {
	// update fields for non-deleted users only
	return r.db.Model(&domain.User{}).Where("id = ? AND status != ?", userID, "deleted").Updates(fields).Error
}