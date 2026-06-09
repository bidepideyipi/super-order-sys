package main

import (
	"fmt"
	"log"
	"super-order-web/internal/config"
	"super-order-web/internal/handler"
	"super-order-web/internal/routes"
	"super-order-web/internal/service"
	"super-order-web/pkg/database"
	"super-order-web/pkg/oss"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	configPath := config.GetDefaultConfigPath()
	if err := config.Load(configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	cfg := config.Get()

	// 初始化数据库
	if err := database.Initialize(&cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化OSS（可选）
	if cfg.OSS.Endpoint != "" {
		if err := oss.Initialize(&cfg.OSS); err != nil {
			log.Printf("初始化OSS失败: %v", err)
		}
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化服务层
	db := database.GetDB()
	customerService := service.NewCustomerService(db)
	skuCategoryService := service.NewSKUCategoryService(db)
	skuService := service.NewSKUService(db)
	orderService := service.NewOrderService(db)
	financialTransactionService := service.NewFinancialTransactionService(db)

	// 初始化处理器
	customerHandler := handler.NewCustomerHandler(customerService)
	skuCategoryHandler := handler.NewSKUCategoryHandler(skuCategoryService)
	skuHandler := handler.NewSKUHandler(skuService)
	orderHandler := handler.NewOrderHandler(orderService)
	financialTransactionHandler := handler.NewFinancialTransactionHandler(financialTransactionService)

	// 设置路由
	router := routes.Setup(
		customerHandler,
		skuCategoryHandler,
		skuHandler,
		orderHandler,
		financialTransactionHandler,
	)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
