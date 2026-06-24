package main

import (
	"fmt"
	"log"

	"super-order-web/internal/config"
	"super-order-web/internal/handler"
	"super-order-web/internal/routes"
	"super-order-web/internal/service"
	"super-order-web/pkg/database"
	"super-order-web/pkg/jwt"
	"super-order-web/pkg/oss"

	_ "super-order-web/docs"
	"github.com/gin-gonic/gin"
)

// @title           超级订单管理系统 API
// @version         1.0
// @description     这是一个超级订单管理系统的后端 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8173
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置
	configPath := config.GetDefaultConfigPath()
	if err := config.Load(configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	cfg := config.Get()

	// 初始化JWT
	jwt.Initialize(&cfg.User)

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
	orderItemService := service.NewOrderItemService(db)
	financialTransactionService := service.NewFinancialTransactionService(db)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(&cfg.User)
	customerHandler := handler.NewCustomerHandler(customerService)
	skuCategoryHandler := handler.NewSKUCategoryHandler(skuCategoryService)
	skuHandler := handler.NewSKUHandler(skuService)
	orderHandler := handler.NewOrderHandler(orderService)
	orderItemHandler := handler.NewOrderItemHandler(orderItemService)
	financialHandler := handler.NewFinancialHandler(financialTransactionService)
	commonHandler := handler.NewCommonHandler(skuService)

	// 设置路由
	router := routes.Setup(
		authHandler,
		customerHandler,
		skuCategoryHandler,
		skuHandler,
		orderHandler,
		orderItemHandler,
		financialHandler,
		commonHandler,
	)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
