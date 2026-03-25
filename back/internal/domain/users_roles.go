package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRoles struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID uuid.UUID `gorm:"type:uuid;primaryKey"`

	AssignedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	User *User `gorm:"foreignKey:UserID"`
	Role *Role `gorm:"foreignKey:RoleID"`
}