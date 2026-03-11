package service

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
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
	repo      *repository.Repository
	config    *config.Config
	wechatPay *utils.PayClient
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	s := &Service{
		repo:   repo,
		config: cfg,
	}

	// 初始化微信支付客户端
	if cfg.Wechat.AppID != "" && cfg.Wechat.MchID != "" && cfg.Wechat.PrivateKeyPath != "" {
		log.Printf("[DEBUG] 初始化微信支付, NotifyURL=%s", cfg.Wechat.NotifyURL)
		client, err := utils.NewPayClient(
			cfg.Wechat.AppID,
			cfg.Wechat.MchID,
			cfg.Wechat.APIV3Key,
			cfg.Wechat.SerialNo,
			cfg.Wechat.PrivateKeyPath,
			cfg.Wechat.NotifyURL,
		)
		if err != nil {
			log.Printf("初始化微信支付客户端失败: %v", err)
		} else {
			s.wechatPay = client
			log.Printf("微信支付客户端初始化成功, NotifyURL=%s", cfg.Wechat.NotifyURL)
		}
	}

	return s
}

// GetOpenIDByCode 通过code获取openid
func (s *Service) GetOpenIDByCode(code string) (string, error) {
	if s.config.Wechat.AppID == "" || s.config.Wechat.AppSecret == "" {
		return "", errors.New("微信配置不完整")
	}
	return utils.GetOpenIDByCode(s.config.Wechat.AppID, s.config.Wechat.AppSecret, code)
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

	// 前端已 MD5 加密，直接比对
	passwordMatch := password == user.PasswordHash
	log.Printf("[LOGIN] 输入密码: %s", password)
	log.Printf("[LOGIN] 数据库密码: %s", user.PasswordHash)
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

// ChangePassword 修改密码
// MD5 加密密码
func md5Hash(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Service) ChangePassword(userID uint, oldPassword, newPassword string) error {
	log.Printf("[ChangePassword service] userID=%d, oldPassword=%s, newPassword=%s", userID, oldPassword, newPassword)

	user, err := s.repo.FindAdminByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	log.Printf("[ChangePassword service] 数据库密码: %s", user.PasswordHash)

	// 前端已 MD5 加密，直接比对
	if oldPassword != user.PasswordHash {
		return errors.New("原密码错误")
	}

	// 直接存前端传的新密码（已 MD5 加密）
	user.PasswordHash = newPassword
	log.Printf("[ChangePassword service] 新密码: %s", newPassword)

	err = s.repo.UpdateAdminUser(user)
	if err != nil {
		log.Printf("[ChangePassword service] 更新失败: %v", err)
		return err
	}

	log.Println("[ChangePassword service] 更新成功")
	return nil
}

// === Card Service Methods ===

func (s *Service) QueryCard(keyword string) (*model.SimCard, error) {
	return s.repo.FindCardByKeyword(keyword)
}

func (s *Service) GetCardByID(id uint) (*model.SimCard, error) {
	return s.repo.FindCardByID(id)
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

	// 2. 获取充值价格（使用卡片记录中的价格）
	price := card.LastRechargeAmount
	if price <= 0 {
		return nil, nil, errors.New("该卡片未设置充值金额，请联系管理员")
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

	// 5. 调用微信统一下单API
	var payParams map[string]interface{}

	log.Printf("[DEBUG] wechatPay=%v, AppID=%s", s.wechatPay, s.config.Wechat.AppID)

	if s.wechatPay != nil {
		// 金额单位是分
		amountFen := int(price * 100)

		prepayResp, err := s.wechatPay.CreateJSAPIPrepay(
			tradeNo,
			"物联网卡充值",
			openid,
			amountFen,
		)
		if err != nil {
			// 微信下单失败，删除已创建的充值记录
			s.repo.DeleteRechargeRecord(record.ID)
			return nil, nil, fmt.Errorf("微信下单失败: %w", err)
		}

		// 生成前端支付参数
		payParams, err = s.wechatPay.GeneratePayParams(prepayResp.PrepayID)
		if err != nil {
			// 生成参数失败，删除已创建的充值记录
			s.repo.DeleteRechargeRecord(record.ID)
			return nil, nil, fmt.Errorf("生成支付参数失败: %w", err)
		}
	} else {
		// 模拟模式（微信支付未配置）
		payParams = map[string]interface{}{
			"appId":     s.config.Wechat.AppID,
			"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
			"nonceStr":  fmt.Sprintf("%d", time.Now().UnixNano()),
			"package":   "prepay_id=mock_" + tradeNo,
			"signType":  "RSA",
			"paySign":   "mock_signature",
		}
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

// DeleteRechargeRecord 删除充值记录
func (s *Service) DeleteRechargeRecord(id uint) error {
	return s.repo.DeleteRechargeRecord(id)
}

// HasWechatPay 检查是否配置了微信支付
func (s *Service) HasWechatPay() bool {
	return s.wechatPay != nil
}

// VerifyWechatNotify 验证微信支付回调
func (s *Service) VerifyWechatNotify(headers map[string]string, body string) error {
	if s.wechatPay == nil {
		return nil
	}

	timestamp := headers["Wechatpay-Timestamp"]
	nonce := headers["Wechatpay-Nonce"]
	signature := headers["Wechatpay-Signature"]
	serialNo := headers["Wechatpay-Serial"]

	if timestamp == "" || nonce == "" || signature == "" || serialNo == "" {
		return fmt.Errorf("缺少微信支付回调头信息")
	}

	// 构造验签字符串
	// 需要微信平台公钥来验证签名，这里简化处理
	// 实际应该从微信平台获取公钥
	log.Printf("收到微信支付回调: timestamp=%s, nonce=%s, serial=%s", timestamp, nonce, serialNo)

	return nil
}

// DecryptWechatNotify 解密微信支付回调
func (s *Service) DecryptWechatNotify(ciphertext, nonce, associatedData string) ([]byte, error) {
	if s.wechatPay == nil {
		return nil, fmt.Errorf("微信支付未配置")
	}

	return s.wechatPay.DecryptNotification(ciphertext, nonce, associatedData)
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
