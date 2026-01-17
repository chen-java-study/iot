package router

import (
	"github.com/gin-gonic/gin"
	"iot-card-system/internal/handler"
	"iot-card-system/internal/middleware"
)

func SetupRouter(h *handler.Handler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

	// API v1
	v1 := r.Group("/api/v1")
	{
		// H5端API (无需认证)
		v1.GET("/card/query", h.QueryCard)
		v1.POST("/payment/create", h.CreatePaymentOrder)
		v1.POST("/payment/notify", h.WechatPaymentNotify)
		v1.GET("/payment/status", h.QueryPaymentStatus)

		// 管理端API
		admin := v1.Group("/admin")
		{
			// 登录接口 (无需认证)
			admin.POST("/login", h.AdminLogin)

			// 需要认证的接口
			auth := admin.Group("", middleware.JWTAuth(jwtSecret))
			{
				auth.GET("/statistics", h.GetStatistics)

				// 卡片管理
				auth.GET("/cards", h.ListCards)
				auth.POST("/cards", h.CreateCard)
				auth.PUT("/cards/:id", h.UpdateCard)
				auth.DELETE("/cards/:id", h.DeleteCard)

				// 充值记录
				auth.GET("/recharges", h.ListRechargeRecords)

				// 系统配置
				auth.GET("/config", h.GetConfig)
				auth.POST("/config", h.UpdateConfig)
			}
		}
	}

	return r
}
