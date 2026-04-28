package repository

import (
	"backend-app/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UsersGroupsRepository struct {
	db *gorm.DB
}

func NewUsersGroupsRepository(db *gorm.DB) *UsersGroupsRepository {
	return &UsersGroupsRepository{db: db}
}

func (r *UsersGroupsRepository) AddUserToGroup(userID uuid.UUID, groupID uuid.UUID, role string) error {
	var ug domain.UsersGroups

	// Unscoped to allow revive of soft-deleted membership.
	err := r.db.Unscoped().
		Where("user_id = ? AND group_id = ?", userID, groupID).
		First(&ug).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(&domain.UsersGroups{
				UserID: userID,
				GroupID: groupID,
				Role:   role,
			}).Error
		}
		return err
	}

	ug.Role = role
	ug.DeletedAt = gorm.DeletedAt{}
	ug.UpdatedAt = time.Now()
	return r.db.Save(&ug).Error
}

func (r *UsersGroupsRepository) RemoveUserFromGroup(userID uuid.UUID, groupID uuid.UUID) error {
	return r.db.
		Where("user_id = ? AND group_id = ? AND deleted_at IS NULL", userID, groupID).
		Delete(&domain.UsersGroups{}).Error
}

func (r *UsersGroupsRepository) GetActiveGroupsByUserID(userID uuid.UUID) ([]*domain.Group, error) {
	var memberships []domain.UsersGroups
	err := r.db.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Preload("Group", "deleted_at IS NULL").
		Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	groups := make([]*domain.Group, 0, len(memberships))
	for i := range memberships {
		if memberships[i].Group != nil {
			groups = append(groups, memberships[i].Group)
		}
	}
	return groups, nil
}