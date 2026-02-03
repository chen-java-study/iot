package repository

import (
	"gorm.io/gorm"
	"iot-card-system/internal/model"
	"strconv"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// === Admin User Methods ===

func (r *Repository) FindAdminByUsername(username string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// UpdateAdminUser 更新管理员用户信息
func (r *Repository) UpdateAdminUser(user *model.AdminUser) error {
	return r.db.Save(user).Error
}

// === SimCard Methods ===

func (r *Repository) FindCardByKeyword(keyword string) (*model.SimCard, error) {
	var card model.SimCard
	err := r.db.Where("card_no = ? OR device_no = ?", keyword, keyword).First(&card).Error
	return &card, err
}

func (r *Repository) FindCardByID(id uint) (*model.SimCard, error) {
	var card model.SimCard
	err := r.db.First(&card, id).Error
	return &card, err
}

func (r *Repository) CreateCard(card *model.SimCard) error {
	return r.db.Create(card).Error
}

func (r *Repository) UpdateCard(card *model.SimCard) error {
	return r.db.Save(card).Error
}

func (r *Repository) DeleteCard(id uint) error {
	return r.db.Delete(&model.SimCard{}, id).Error
}

func (r *Repository) ListCards(page, pageSize int, status int16, keyword string) ([]model.SimCard, int64, error) {
	var cards []model.SimCard
	var total int64

	query := r.db.Model(&model.SimCard{})

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if keyword != "" {
		query = query.Where("card_no LIKE ? OR device_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&cards).Error

	return cards, total, err
}

// === RechargeRecord Methods ===

func (r *Repository) CreateRechargeRecord(record *model.RechargeRecord) error {
	return r.db.Create(record).Error
}

func (r *Repository) FindRechargeByTradeNo(tradeNo string) (*model.RechargeRecord, error) {
	var record model.RechargeRecord
	err := r.db.Where("trade_no = ?", tradeNo).First(&record).Error
	return &record, err
}

func (r *Repository) UpdateRechargeRecord(record *model.RechargeRecord) error {
	return r.db.Save(record).Error
}

func (r *Repository) ListRechargeRecords(page, pageSize int, status int16, keyword, startDate, endDate string) ([]model.RechargeRecord, int64, float64, error) {
	var records []model.RechargeRecord
	var total int64
	var totalAmount float64

	query := r.db.Model(&model.RechargeRecord{})

	if status >= 0 {
		query = query.Where("payment_status = ?", status)
	}

	if keyword != "" {
		query = query.Where("card_no LIKE ? OR device_no LIKE ? OR trade_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	query.Count(&total)

	// 计算总金额 - 使用相同的查询条件构建新查询
	amountQuery := r.db.Model(&model.RechargeRecord{})
	if status >= 0 {
		amountQuery = amountQuery.Where("payment_status = ?", status)
	}
	if keyword != "" {
		amountQuery = amountQuery.Where("card_no LIKE ? OR device_no LIKE ? OR trade_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if startDate != "" {
		amountQuery = amountQuery.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		amountQuery = amountQuery.Where("DATE(created_at) <= ?", endDate)
	}
	amountQuery.Select("COALESCE(SUM(recharge_amount), 0)").Scan(&totalAmount)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error

	return records, total, totalAmount, err
}

// === SystemConfig Methods ===

func (r *Repository) GetConfig(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&config).Error
	return &config, err
}

func (r *Repository) GetAllConfigs() (map[string]string, error) {
	var configs []model.SystemConfig
	if err := r.db.Find(&configs).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.ConfigValue
	}
	return result, nil
}

func (r *Repository) UpdateConfig(key, value string) error {
	return r.db.Model(&model.SystemConfig{}).Where("config_key = ?", key).Update("config_value", value).Error
}

func (r *Repository) GetConfigFloat(key string) (float64, error) {
	config, err := r.GetConfig(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(config.ConfigValue, 64)
}

// === Statistics Methods ===

func (r *Repository) GetStatistics() (*model.Statistics, error) {
	var stats model.Statistics
	err := r.db.Raw("SELECT * FROM v_statistics").Scan(&stats).Error
	return &stats, err
}

// === Transaction Support ===

func (r *Repository) Transaction(fn func(*Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{db: tx}
		return fn(txRepo)
	})
}
