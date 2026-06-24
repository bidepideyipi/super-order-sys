package routes

import (
	"super-order-web/internal/handler"
	"super-order-web/internal/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// Setup 设置路由
func Setup(
	authHandler *handler.AuthHandler,
	customerHandler *handler.CustomerHandler,
	skuCategoryHandler *handler.SKUCategoryHandler,
	skuHandler *handler.SKUHandler,
	orderHandler *handler.OrderHandler,
	orderItemHandler *handler.OrderItemHandler,
	financialHandler *handler.FinancialHandler,
	commonHandler *handler.CommonHandler,
) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())

	// API路由组
	api := router.Group("/api")
	{
		// ========== 认证路由 ==========
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/check", middleware.Auth(), authHandler.Check)
		}

		// ========== SKU 路由 ==========
		sku := api.Group("/sku")
		sku.Use(middleware.Auth())
		{
			sku.GET("/list", skuHandler.ListAll)
			sku.GET("/list-paginated", skuHandler.ListPaginated)
			sku.GET("/search", skuHandler.Search)
			sku.GET("/search-paginated", skuHandler.SearchPaginated)
			sku.GET("/search-with-category", skuHandler.SearchWithCategory)
			sku.GET("/:id", skuHandler.Get)
			sku.POST("", skuHandler.Create)
			sku.PUT("/:id", skuHandler.Update)
			sku.DELETE("/:id", skuHandler.Delete)
		}

		// ========== 分类路由 ==========
		category := api.Group("/category")
		category.Use(middleware.Auth())
		{
			category.GET("/list", skuCategoryHandler.List)
			category.GET("/:id", skuCategoryHandler.Get)
			category.POST("", skuCategoryHandler.Create)
			category.PUT("/:id", skuCategoryHandler.Update)
			category.DELETE("/:id", skuCategoryHandler.Delete)
		}

		// ========== 客户路由 ==========
		customer := api.Group("/customer")
		customer.Use(middleware.Auth())
		{
			customer.GET("/list", customerHandler.List)
			customer.GET("/:id", customerHandler.Get)
			customer.POST("", customerHandler.Create)
			customer.PUT("/:id", customerHandler.Update)
			customer.DELETE("/:id", customerHandler.Delete)
		}

		// ========== 订单路由 ==========
		order := api.Group("/order")
		order.Use(middleware.Auth())
		{
			order.GET("/list", orderHandler.List)
			order.GET("/processing", orderHandler.GetProcessingOrders)
			order.GET("/unsettled", orderHandler.GetUnsettledOrders)
			order.GET("/:id", orderHandler.Get)
			order.GET("/:id/items", orderHandler.GetItems)
			order.POST("", orderHandler.Create)
			order.PUT("/:id", orderHandler.Update)
			order.DELETE("/:id", orderHandler.Delete)
		}

		// ========== 订单明细路由 ==========
		orderItem := api.Group("/order-item")
		orderItem.Use(middleware.Auth())
		{
			orderItem.GET("/:id", orderItemHandler.Get)
			orderItem.POST("", orderItemHandler.Create)
			orderItem.PUT("/:id", orderItemHandler.Update)
			orderItem.DELETE("/:id", orderItemHandler.Delete)
		}

		// ========== 财务路由 ==========
		financial := api.Group("/financial")
		financial.Use(middleware.Auth())
		{
			financial.GET("/list", financialHandler.List)
			financial.GET("/balance", financialHandler.GetBalance)
			financial.GET("/:id", financialHandler.Get)
			financial.POST("", financialHandler.Create)
			financial.PUT("/:id", financialHandler.Update)
			financial.DELETE("/:id", financialHandler.Delete)
		}

		// ========== 公共路由 ==========
		common := api.Group("/common")
		{
			common.GET("/image/:skuCode", commonHandler.GetImage)
		}
	}

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
