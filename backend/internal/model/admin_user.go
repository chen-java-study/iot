package model

import (
	"time"
)

type AdminUser struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Username    string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	RealName    string     `gorm:"size:100" json:"real_name"`
	Phone       string     `gorm:"size:20" json:"phone"`
	Email       string     `gorm:"size:100" json:"email"`
	Status      int16      `gorm:"default:1" json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}
