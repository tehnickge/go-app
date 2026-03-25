package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"size:50;not null;unique"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	UserRoles []*UsersRoles `gorm:"foreignKey:RoleID"`
}