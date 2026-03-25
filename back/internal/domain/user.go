package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Login     string         `gorm:"size:100;not null;unique"`
	Email     string         `gorm:"size:255;not null;unique"`
	Password  string         `gorm:"size:255;not null"`
	Status    string         `gorm:"type:user_status;default:'active'"`
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserGroups []*UsersGroups `gorm:"foreignKey:UserID"`
	UserRoles  []*UsersRoles  `gorm:"foreignKey:UserID"`
}