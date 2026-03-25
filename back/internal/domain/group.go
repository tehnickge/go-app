package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type Group struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string         `gorm:"size:255;not null"`
	Description string
	OwnerID     *uuid.UUID
	Owner       *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	UserGroups []*UsersGroups `gorm:"foreignKey:GroupID"`
}