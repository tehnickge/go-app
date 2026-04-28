package repository

import (
	"backend-app/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UsersRolesRepository struct {
	db *gorm.DB
}

func NewUsersRolesRepository(db *gorm.DB) *UsersRolesRepository {
	return &UsersRolesRepository{db: db}
}

func (r *UsersRolesRepository) AssignRoleToUser(userID uuid.UUID, roleID uuid.UUID) error {
	var ur domain.UsersRoles

	// Unscoped to be able to "revive" soft-deleted assignments.
	err := r.db.Unscoped().
		Where("user_id = ? AND role_id = ?", userID, roleID).
		First(&ur).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(&domain.UsersRoles{
				UserID: userID,
				RoleID: roleID,
			}).Error
		}
		return err
	}

	ur.DeletedAt = gorm.DeletedAt{}
	ur.UpdatedAt = time.Now()
	return r.db.Save(&ur).Error
}

func (r *UsersRolesRepository) RemoveRoleFromUser(userID uuid.UUID, roleID uuid.UUID) error {
	return r.db.
		Where("user_id = ? AND role_id = ? AND deleted_at IS NULL", userID, roleID).
		Delete(&domain.UsersRoles{}).Error
}

func (r *UsersRolesRepository) GetActiveRolesByUserID(userID uuid.UUID) ([]*domain.Role, error) {
	var assignments []domain.UsersRoles
	err := r.db.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Preload("Role", "deleted_at IS NULL").
		Find(&assignments).Error
	if err != nil {
		return nil, err
	}

	roles := make([]*domain.Role, 0, len(assignments))
	for i := range assignments {
		if assignments[i].Role != nil {
			roles = append(roles, assignments[i].Role)
		}
	}
	return roles, nil
}

func (r *UsersRolesRepository) HasAnyActiveAssignmentByRoleID(roleID uuid.UUID) (bool, error) {
	var cnt int64
	if err := r.db.
		Model(&domain.UsersRoles{}).
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}