package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"iot-card-system/internal/model"
	"iot-card-system/internal/service"
	"iot-card-system/internal/utils"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// === H5端接口 ===

// QueryCard 查询卡片
func (h *Handler) QueryCard(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.BadRequest(c, "请输入卡号或设备号")
		return
	}

	card, err := h.service.QueryCard(keyword)
	if err != nil {
		utils.BadRequest(c, "未找到卡号「"+keyword+"」，请检查输入是否正确")
		return
	}

	utils.Success(c, gin.H{
		"id":                  card.ID,
		"card_no":             card.CardNo,
		"device_no":           card.DeviceNo,
		"start_date":          card.StartDate.Format("2006-01-02"),
		"expire_date":         card.ExpireDate.Format("2006-01-02"),
		"status":              card.Status,
		"status_text":         card.StatusText(),
		"operator":            card.Operator,
		"days_remaining":     card.DaysRemaining(),
		"last_recharge_amount": card.LastRechargeAmount,
	})
}

// GetOpenIDByCode 通过微信授权code获取openid
func (h *Handler) GetOpenIDByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.BadRequest(c, "缺少code参数")
		return
	}

	openid, err := h.service.GetOpenIDByCode(code)
	if err != nil {
		utils.Error(c, 500, "获取openid失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"openid": openid})
}

// CreatePaymentOrder 创建充值订单
func (h *Handler) CreatePaymentOrder(c *gin.Context) {
	var req struct {
		CardNo string `json:"card_no" binding:"required"`
		Openid string `json:"openid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	record, payParams, err := h.service.CreateRechargeOrder(req.CardNo, req.Openid, ipAddress, userAgent)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"trade_no":   record.TradeNo,
		"amount":     record.RechargeAmount,
		"pay_params": payParams,
	})
}

// WechatPaymentNotify 微信支付回调
func (h *Handler) WechatPaymentNotify(c *gin.Context) {
	// 获取微信回调请求头
	headers := map[string]string{
		"Wechatpay-Timestamp": c.GetHeader("Wechatpay-Timestamp"),
		"Wechatpay-Nonce":     c.GetHeader("Wechatpay-Nonce"),
		"Wechatpay-Serial":    c.GetHeader("Wechatpay-Serial"),
		"Wechatpay-Signature": c.GetHeader("Wechatpay-Signature"),
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"code": "FAIL", "message": "读取请求体失败"})
		return
	}

	// 验证签名（如果配置了微信支付）
	if h.service.HasWechatPay() {
		err := h.service.VerifyWechatNotify(headers, string(body))
		if err != nil {
			log.Printf("微信支付回调验签失败: %v", err)
			c.JSON(400, gin.H{"code": "FAIL", "message": "验签失败"})
			return
		}
	}

	// 解析回调数据
	var notifyData struct {
		ID           string `json:"id"`
		CreateTime   string `json:"create_time"`
		ResourceType string `json:"resource_type"`
		EventType    string `json:"event_type"`
		Resource     struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			Nonce          string `json:"nonce"`
			AssociatedData string `json:"associated_data"`
		} `json:"resource"`
		Summary string `json:"summary"`
	}

	if err := json.Unmarshal(body, &notifyData); err != nil {
		c.JSON(400, gin.H{"code": "FAIL", "message": "解析回调数据失败"})
		return
	}

	// 只处理支付成功通知
	if notifyData.EventType != "TRANSACTION.SUCCESS" {
		c.JSON(200, gin.H{"code": "SUCCESS", "message": "ok"})
		return
	}

	// 解密通知数据
	plaintext, err := h.service.DecryptWechatNotify(
		notifyData.Resource.Ciphertext,
		notifyData.Resource.Nonce,
		notifyData.Resource.AssociatedData,
	)
	if err != nil {
		log.Printf("解密微信支付回调失败: %v", err)
		c.JSON(500, gin.H{"code": "FAIL", "message": "解密失败"})
		return
	}

	// 解析明文
	var paymentResult struct {
		OutTradeNo     string `json:"out_trade_no"`
		TransactionID  string `json:"transaction_id"`
		TradeState     string `json:"trade_state"`
		TradeStateDesc string `json:"trade_state_desc"`
		Amount         struct {
			Total         int    `json:"total"`
			PayerTotal    int    `json:"payer_total"`
			Currency      string `json:"currency"`
			PayerCurrency string `json:"payer_currency"`
		} `json:"amount"`
	}

	if err := json.Unmarshal(plaintext, &paymentResult); err != nil {
		log.Printf("解析支付结果失败: %v", err)
		c.JSON(500, gin.H{"code": "FAIL", "message": "解析支付结果失败"})
		return
	}

	// 处理支付成功
	if paymentResult.TradeState == "SUCCESS" {
		paidAt, _ := time.Parse("2006-01-02T15:04:05+08:00", notifyData.CreateTime)
		err := h.service.HandlePaymentNotify(paymentResult.TransactionID, paymentResult.OutTradeNo, paidAt)
		if err != nil {
			log.Printf("处理支付通知失败: %v", err)
			c.JSON(500, gin.H{"code": "FAIL", "message": "处理失败"})
			return
		}
	}

	c.JSON(200, gin.H{"code": "SUCCESS", "message": "ok"})
}

// QueryPaymentStatus 查询订单状态
func (h *Handler) QueryPaymentStatus(c *gin.Context) {
	tradeNo := c.Query("trade_no")
	if tradeNo == "" {
		utils.BadRequest(c, "参数错误")
		return
	}

	record, err := h.service.QueryPaymentStatus(tradeNo)
	if err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	utils.Success(c, gin.H{
		"trade_no":            record.TradeNo,
		"payment_status":      record.PaymentStatus,
		"payment_status_text": record.PaymentStatusText(),
		"paid_at":             record.PaidAt,
	})
}

// === 管理端接口 ===

// AdminLogin 管理员登录
func (h *Handler) AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	token, user, err := h.service.AdminLogin(req.Username, req.Password)
	if err != nil {
		fmt.Println("AdminLogin error: ", err)
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token":  token,
		"expire": time.Now().Add(24 * time.Hour).Unix(),
		"user_info": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
		},
	})
}

// ChangePassword 修改密码
func (h *Handler) ChangePassword(c *gin.Context) {
	log.Println("[ChangePassword] 收到修改密码请求")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ChangePassword] 参数错误: %v", err)
		utils.BadRequest(c, "参数错误")
		return
	}

	log.Printf("[ChangePassword] 旧密码: %s, 新密码: %s", req.OldPassword, req.NewPassword)

	// 从 JWT 中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("[ChangePassword] 未登录")
		utils.Unauthorized(c, "未登录")
		return
	}

	log.Printf("[ChangePassword] 用户ID: %v", userID)

	err := h.service.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		log.Printf("[ChangePassword] 修改失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	log.Println("[ChangePassword] 修改成功")
	utils.Success(c, "密码修改成功")
}

// GetStatistics 获取统计数据
func (h *Handler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

// ListCards 卡片列表
func (h *Handler) ListCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	keyword := c.Query("keyword")

	cards, total, err := h.service.ListCards(page, pageSize, int16(status), keyword)
	if err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	// 格式化日期
	var result []gin.H
	for _, card := range cards {
		lastRechargeTimeStr := ""
		if card.LastRechargeTime != nil {
			lastRechargeTimeStr = card.LastRechargeTime.Format("2006-01-02 15:04:05")
		}
		result = append(result, gin.H{
			"id":                    card.ID,
			"card_no":               card.CardNo,
			"device_no":             card.DeviceNo,
			"start_date":            card.StartDate.Format("2006-01-02"),
			"expire_date":           card.ExpireDate.Format("2006-01-02"),
			"status":                card.Status,
			"status_text":           card.StatusText(),
			"operator":              card.Operator,
			"package_type":          card.PackageType,
			"total_recharge_count":  card.TotalRechargeCount,
			"total_recharge_amount": card.TotalRechargeAmount,
			"last_recharge_time":    lastRechargeTimeStr,
			"last_recharge_amount":  card.LastRechargeAmount,
			"remark":                card.Remark,
			"created_at":            card.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.Success(c, gin.H{
		"list":      result,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateCard 创建卡片
func (h *Handler) CreateCard(c *gin.Context) {
	var req struct {
		CardNo             string  `json:"card_no"`
		DeviceNo           string  `json:"device_no"`
		Operator           string  `json:"operator"`
		PackageType        string  `json:"package_type"`
		StartDate          string  `json:"start_date"`
		ExpireDate         string  `json:"expire_date"`
		LastRechargeAmount float64 `json:"last_recharge_amount"`
		Remark             string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.BadRequest(c, "开始日期格式错误")
		return
	}
	expireDate, err := time.Parse("2006-01-02", req.ExpireDate)
	if err != nil {
		utils.BadRequest(c, "到期日期格式错误")
		return
	}

	card := model.SimCard{
		CardNo:             req.CardNo,
		DeviceNo:           req.DeviceNo,
		Operator:           req.Operator,
		PackageType:        req.PackageType,
		StartDate:          startDate,
		ExpireDate:         expireDate,
		LastRechargeAmount: req.LastRechargeAmount,
		Remark:             req.Remark,
	}

	// 有充值金额时，自动记录当前时间为充值时间
	if req.LastRechargeAmount > 0 {
		now := time.Now()
		card.LastRechargeTime = &now
	}

	if err := h.service.CreateCard(&card); err != nil {
		utils.BadRequest(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, card)
}

// UpdateCard 更新卡片
func (h *Handler) UpdateCard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		CardNo             string  `json:"card_no"`
		DeviceNo           string  `json:"device_no"`
		Operator           string  `json:"operator"`
		PackageType        string  `json:"package_type"`
		StartDate          string  `json:"start_date"`
		ExpireDate         string  `json:"expire_date"`
		LastRechargeAmount float64 `json:"last_recharge_amount"`
		Remark             string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 先查询原有卡片，保留关键字段
	oldCard, err := h.service.GetCardByID(uint(id))
	if err != nil {
		utils.BadRequest(c, "卡片不存在")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.BadRequest(c, "开始日期格式错误")
		return
	}
	expireDate, err := time.Parse("2006-01-02", req.ExpireDate)
	if err != nil {
		utils.BadRequest(c, "到期日期格式错误")
		return
	}

	card := model.SimCard{
		ID:                  uint(id),
		CardNo:              req.CardNo,
		DeviceNo:            req.DeviceNo,
		Operator:            req.Operator,
		PackageType:         req.PackageType,
		StartDate:           startDate,
		ExpireDate:          expireDate,
		LastRechargeAmount:  req.LastRechargeAmount,
		Remark:              req.Remark,
		Status:              oldCard.Status,
		TotalRechargeCount:  oldCard.TotalRechargeCount,
		TotalRechargeAmount: oldCard.TotalRechargeAmount,
		LastRechargeTime:    oldCard.LastRechargeTime,
	}

	// 有充值金额时，自动记录当前时间为充值时间
	if req.LastRechargeAmount > 0 {
		now := time.Now()
		card.LastRechargeTime = &now
	}

	if err := h.service.UpdateCard(&card); err != nil {
		utils.BadRequest(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, card)
}

// DeleteCard 删除卡片
func (h *Handler) DeleteCard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteCard(uint(id)); err != nil {
		utils.BadRequest(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}

// ListRechargeRecords 充值记录列表
func (h *Handler) ListRechargeRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("payment_status", "-1"))
	keyword := c.Query("keyword")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	records, total, totalAmount, err := h.service.ListRechargeRecords(page, pageSize, int16(status), keyword, startDate, endDate)
	if err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	var result []gin.H
	for _, record := range records {
		paidAtStr := ""
		if record.PaidAt != nil {
			paidAtStr = record.PaidAt.Format("2006-01-02 15:04:05")
		}
		result = append(result, gin.H{
			"id":                  record.ID,
			"card_no":             record.CardNo,
			"device_no":           record.DeviceNo,
			"recharge_amount":     record.RechargeAmount,
			"trade_no":            record.TradeNo,
			"payment_status":      record.PaymentStatus,
			"payment_status_text": record.PaymentStatusText(),
			"paid_at":             paidAtStr,
			"created_at":          record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.Success(c, gin.H{
		"list":         result,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"total_amount": totalAmount,
	})
}

// 允许通过接口读写的配置项白名单（微信支付相关敏感配置只在 config.yaml 里维护）
var allowedConfigKeys = map[string]bool{
	"alert_days": true,
}

// GetConfig 获取系统配置
func (h *Handler) GetConfig(c *gin.Context) {
	configs, err := h.service.GetAllConfigs()
	if err != nil {
		utils.InternalError(c, "获取配置失败")
		return
	}

	// 只返回白名单内的配置，屏蔽微信支付等敏感信息
	safeConfigs := make(map[string]string)
	for k, v := range configs {
		if allowedConfigKeys[k] {
			safeConfigs[k] = v
		}
	}

	utils.Success(c, safeConfigs)
}

// UpdateConfig 更新系统配置
func (h *Handler) UpdateConfig(c *gin.Context) {
	var configs map[string]string
	if err := c.ShouldBindJSON(&configs); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 只允许更新白名单内的配置
	safeConfigs := make(map[string]string)
	for k, v := range configs {
		if allowedConfigKeys[k] {
			safeConfigs[k] = v
		}
	}

	if err := h.service.UpdateConfigs(safeConfigs); err != nil {
		utils.InternalError(c, "更新配置失败")
		return
	}

	utils.Success(c, nil)
}
