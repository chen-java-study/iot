package main

import (
	"fmt"
	"iot-card-system/internal/config"
	"iot-card-system/internal/handler"
	"iot-card-system/internal/repository"
	"iot-card-system/internal/router"
	"iot-card-system/internal/service"
	"iot-card-system/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 环境变量文件
	if err := godotenv.Load(); err != nil {
		log.Printf("加载 .env 文件失败: %v", err)
	}

	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := database.NewPostgresDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 初始化Repository
	repo := repository.NewRepository(db)

	// 初始化Service
	svc := service.NewService(repo, cfg)

	// 初始化Handler
	h := handler.NewHandler(svc)

	// 设置路由
	gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter(h, cfg.JWT.SecretKey)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在端口 %d", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
