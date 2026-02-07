package handler

import (
	"fmt"
	"iot-card-system/internal/model"
	"iot-card-system/internal/service"
	"iot-card-system/internal/utils"
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
		utils.NotFound(c, "卡片不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":             card.ID,
		"card_no":        card.CardNo,
		"device_no":      card.DeviceNo,
		"start_date":     card.StartDate.Format("2006-01-02"),
		"expire_date":    card.ExpireDate.Format("2006-01-02"),
		"status":         card.Status,
		"status_text":    card.StatusText(),
		"operator":       card.Operator,
		"days_remaining": card.DaysRemaining(),
	})
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
	// 简化版实现，实际需要验证签名并解密
	var req struct {
		TransactionID string `json:"transaction_id"`
		OutTradeNo    string `json:"out_trade_no"`
		TradeState    string `json:"trade_state"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": "FAIL", "message": "参数错误"})
		return
	}

	if req.TradeState == "SUCCESS" {
		err := h.service.HandlePaymentNotify(req.TransactionID, req.OutTradeNo, time.Now())
		if err != nil {
			c.JSON(500, gin.H{"code": "FAIL", "message": "处理失败"})
			return
		}
	}

	c.JSON(200, gin.H{"code": "SUCCESS", "message": ""})
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
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
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
		LastRechargeTime   string  `json:"last_recharge_time"`
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

	if req.LastRechargeTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.LastRechargeTime)
		if err == nil {
			card.LastRechargeTime = &t
		}
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
		LastRechargeTime   string  `json:"last_recharge_time"`
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
		ID:                 uint(id),
		CardNo:             req.CardNo,
		DeviceNo:           req.DeviceNo,
		Operator:           req.Operator,
		PackageType:        req.PackageType,
		StartDate:          startDate,
		ExpireDate:         expireDate,
		LastRechargeAmount: req.LastRechargeAmount,
		Remark:             req.Remark,
	}

	if req.LastRechargeTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.LastRechargeTime)
		if err == nil {
			card.LastRechargeTime = &t
		}
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

// GetConfig 获取系统配置
func (h *Handler) GetConfig(c *gin.Context) {
	configs, err := h.service.GetAllConfigs()
	if err != nil {
		utils.InternalError(c, "获取配置失败")
		return
	}

	utils.Success(c, configs)
}

// UpdateConfig 更新系统配置
func (h *Handler) UpdateConfig(c *gin.Context) {
	var configs map[string]string
	if err := c.ShouldBindJSON(&configs); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateConfigs(configs); err != nil {
		utils.InternalError(c, "更新配置失败")
		return
	}

	utils.Success(c, nil)
}
