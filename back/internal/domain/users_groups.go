package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersGroups struct {
	UserID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Role      string    `gorm:"size:50;default:'member'"` // роль в группе
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User  *User  `gorm:"foreignKey:UserID"`
	Group *Group `gorm:"foreignKey:GroupID"`
}