package model

import (
	"time"
)

type RechargeRecord struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CardID          uint       `gorm:"not null;index" json:"card_id"`
	CardNo          string     `gorm:"size:50;not null;index" json:"card_no"`
	DeviceNo        string     `gorm:"size:50;not null" json:"device_no"`
	RechargeAmount  float64    `gorm:"type:decimal(10,2);not null" json:"recharge_amount"`
	RechargeYears   int        `gorm:"default:1" json:"recharge_years"`
	OldExpireDate   time.Time  `gorm:"type:date;not null" json:"old_expire_date"`
	NewExpireDate   time.Time  `gorm:"type:date;not null" json:"new_expire_date"`
	PaymentMethod   string     `gorm:"size:20;default:'wechat'" json:"payment_method"`
	TradeNo         string     `gorm:"uniqueIndex;size:100" json:"trade_no"`
	TransactionID   string     `gorm:"size:100" json:"transaction_id"`
	PaymentStatus   int16      `gorm:"default:0;index" json:"payment_status"`
	PaidAt          *time.Time `gorm:"index" json:"paid_at"`
	Openid          string     `gorm:"size:100" json:"openid"`
	IPAddress       string     `gorm:"size:50" json:"ip_address"`
	UserAgent       string     `gorm:"type:text" json:"user_agent"`
	CreatedAt       time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (RechargeRecord) TableName() string {
	return "recharge_records"
}

// PaymentStatusText 返回支付状态文本
func (r *RechargeRecord) PaymentStatusText() string {
	switch r.PaymentStatus {
	case 0:
		return "待支付"
	case 1:
		return "已支付"
	case 2:
		return "已退款"
	case 3:
		return "支付失败"
	default:
		return "未知"
	}
}
