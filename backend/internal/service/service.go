package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"iot-card-system/internal/config"
	"iot-card-system/internal/model"
	"iot-card-system/internal/repository"
	"iot-card-system/internal/utils"
	"log"
	"math/big"
	"time"
)

type Service struct {
	repo   *repository.Repository
	config *config.Config
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: cfg,
	}
}

// === Admin Service Methods ===

func (s *Service) AdminLogin(username, password string) (string, *model.AdminUser, error) {
	log.Printf("[LOGIN] 尝试登录 - username: %s", username)

	user, err := s.repo.FindAdminByUsername(username)
	if err != nil {
		log.Printf("[LOGIN] 用户不存在 - username: %s, err: %v", username, err)
		return "", nil, errors.New("用户名或密码错误")
	}

	log.Printf("[LOGIN] 找到用户 - id: %d, status: %d", user.ID, user.Status)
	log.Printf("[LOGIN] 数据库中的密码hash: %s", user.PasswordHash)
	log.Printf("[LOGIN] hash长度: %d", len(user.PasswordHash))
	log.Printf("[LOGIN] 输入的密码: %s", password)
	log.Printf("[LOGIN] 密码长度: %d", len(password))

	if user.Status != 1 {
		log.Printf("[LOGIN] 账号已禁用 - status: %d", user.Status)
		return "", nil, errors.New("账号已被禁用")
	}

	// 直接比对明文密码
	passwordMatch := password == user.PasswordHash
	log.Printf("[LOGIN] 密码验证结果: %v", passwordMatch)

	if !passwordMatch {
		log.Printf("[LOGIN] 密码错误")
		return "", nil, errors.New("用户名或密码错误")
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, s.config.JWT.SecretKey, s.config.JWT.ExpireHours)
	if err != nil {
		return "", nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.repo.UpdateAdminUser(user)

	return token, user, nil
}

// === Card Service Methods ===

func (s *Service) QueryCard(keyword string) (*model.SimCard, error) {
	return s.repo.FindCardByKeyword(keyword)
}

func (s *Service) CreateCard(card *model.SimCard) error {
	return s.repo.CreateCard(card)
}

func (s *Service) UpdateCard(card *model.SimCard) error {
	return s.repo.UpdateCard(card)
}

func (s *Service) DeleteCard(id uint) error {
	return s.repo.DeleteCard(id)
}

func (s *Service) ListCards(page, pageSize int, status int16, keyword string) ([]model.SimCard, int64, error) {
	return s.repo.ListCards(page, pageSize, status, keyword)
}

// === Payment Service Methods ===

// GenerateTradeNo 生成订单号
func (s *Service) GenerateTradeNo() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(100000))
	return fmt.Sprintf("R%s%05d", time.Now().Format("20060102150405"), n)
}

// CreateRechargeOrder 创建充值订单
func (s *Service) CreateRechargeOrder(cardNo, openid, ipAddress, userAgent string) (*model.RechargeRecord, map[string]interface{}, error) {
	// 1. 查询卡片
	card, err := s.repo.FindCardByKeyword(cardNo)
	if err != nil {
		return nil, nil, errors.New("卡片不存在")
	}

	// 2. 获取充值价格
	price, err := s.repo.GetConfigFloat("recharge_price")
	if err != nil {
		price = 100.00 // 默认价格
	}

	// 3. 计算新的到期日期
	newExpireDate := card.ExpireDate.AddDate(1, 0, 0)
	if card.ExpireDate.Before(time.Now()) {
		newExpireDate = time.Now().AddDate(1, 0, 0)
	}

	// 4. 创建充值记录
	tradeNo := s.GenerateTradeNo()
	record := &model.RechargeRecord{
		CardID:         card.ID,
		CardNo:         card.CardNo,
		DeviceNo:       card.DeviceNo,
		RechargeAmount: price,
		RechargeYears:  1,
		OldExpireDate:  card.ExpireDate,
		NewExpireDate:  newExpireDate,
		TradeNo:        tradeNo,
		PaymentStatus:  0,
		Openid:         openid,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
	}

	if err := s.repo.CreateRechargeRecord(record); err != nil {
		return nil, nil, err
	}

	// 5. 生成微信支付参数 (简化版，实际需要调用微信API)
	payParams := map[string]interface{}{
		"appId":     s.config.Wechat.AppID,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  fmt.Sprintf("%d", time.Now().UnixNano()),
		"package":   "prepay_id=wx_mock_prepay_id",
		"signType":  "RSA",
		"paySign":   "mock_signature",
	}

	return record, payParams, nil
}

// HandlePaymentNotify 处理微信支付回调
func (s *Service) HandlePaymentNotify(transactionID, tradeNo string, paidAt time.Time) error {
	return s.repo.Transaction(func(repo *repository.Repository) error {
		// 1. 查询充值记录
		record, err := repo.FindRechargeByTradeNo(tradeNo)
		if err != nil {
			return err
		}

		if record.PaymentStatus == 1 {
			return nil // 已处理过
		}

		// 2. 更新充值记录
		record.TransactionID = transactionID
		record.PaymentStatus = 1
		record.PaidAt = &paidAt

		if err := repo.UpdateRechargeRecord(record); err != nil {
			return err
		}

		// 3. 更新卡片到期时间
		card, err := repo.FindCardByID(record.CardID)
		if err != nil {
			return err
		}

		card.ExpireDate = record.NewExpireDate
		card.TotalRechargeCount++
		card.TotalRechargeAmount += record.RechargeAmount
		card.LastRechargeTime = &paidAt
		card.LastRechargeAmount = record.RechargeAmount

		return repo.UpdateCard(card)
	})
}

// QueryPaymentStatus 查询订单状态
func (s *Service) QueryPaymentStatus(tradeNo string) (*model.RechargeRecord, error) {
	return s.repo.FindRechargeByTradeNo(tradeNo)
}

// === Recharge Record Methods ===

func (s *Service) ListRechargeRecords(page, pageSize int, status int16, keyword, startDate, endDate string) ([]model.RechargeRecord, int64, float64, error) {
	return s.repo.ListRechargeRecords(page, pageSize, status, keyword, startDate, endDate)
}

// === Config Methods ===

func (s *Service) GetAllConfigs() (map[string]string, error) {
	return s.repo.GetAllConfigs()
}

func (s *Service) UpdateConfigs(configs map[string]string) error {
	for key, value := range configs {
		if err := s.repo.UpdateConfig(key, value); err != nil {
			return err
		}
	}
	return nil
}

// === Statistics Methods ===

func (s *Service) GetStatistics() (*model.Statistics, error) {
	return s.repo.GetStatistics()
}
