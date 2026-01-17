package model

import (
	"time"
)

type SimCard struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	CardNo              string    `gorm:"uniqueIndex;size:50;not null" json:"card_no"`
	DeviceNo            string    `gorm:"uniqueIndex;size:50;not null" json:"device_no"`
	StartDate           time.Time `gorm:"type:date;not null" json:"start_date"`
	ExpireDate          time.Time `gorm:"type:date;not null" json:"expire_date"`
	Status              int16     `gorm:"default:1" json:"status"`
	Operator            string    `gorm:"size:20" json:"operator"`
	PackageType         string    `gorm:"size:50" json:"package_type"`
	TotalRechargeCount  int       `gorm:"default:0" json:"total_recharge_count"`
	TotalRechargeAmount float64   `gorm:"type:decimal(10,2);default:0" json:"total_recharge_amount"`
	Remark              string    `gorm:"type:text" json:"remark"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (SimCard) TableName() string {
	return "sim_cards"
}

// StatusText 返回状态文本
func (s *SimCard) StatusText() string {
	switch s.Status {
	case 1:
		return "正常"
	case 2:
		return "即将到期"
	case 3:
		return "已过期"
	default:
		return "禁用"
	}
}

// DaysRemaining 计算剩余天数
func (s *SimCard) DaysRemaining() int {
	days := int(time.Until(s.ExpireDate).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}
