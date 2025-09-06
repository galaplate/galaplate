package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string         `gorm:"size:50;not null" json:"username"`
	Email       string         `gorm:"size:100;uniqueIndex" json:"email"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Description string         `gorm:"size:255" json:"description"`
	Status      bool           `gorm:"default:false" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
