package model

import (
	"time"
)

type SystemConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ConfigKey   string    `gorm:"uniqueIndex;size:50;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text;not null" json:"config_value"`
	ConfigType  string    `gorm:"size:20;default:'string'" json:"config_type"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_config"
}

type Statistics struct {
	TotalCards          int64   `json:"total_cards"`
	ActiveCards         int64   `json:"active_cards"`
	ExpiringCards       int64   `json:"expiring_cards"`
	ExpiredCards        int64   `json:"expired_cards"`
	TotalRechargeAmount float64 `json:"total_recharge_amount"`
	TotalRechargeCount  int64   `json:"total_recharge_count"`
	TodayAmount         float64 `json:"today_amount"`
	MonthAmount         float64 `json:"month_amount"`
}
