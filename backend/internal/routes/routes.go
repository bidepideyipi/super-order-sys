package routes

import (
	"super-order-web/internal/handler"
	"super-order-web/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Setup 设置路由
func Setup(
	customerHandler *handler.CustomerHandler,
	skuCategoryHandler *handler.SKUCategoryHandler,
	skuHandler *handler.SKUHandler,
	orderHandler *handler.OrderHandler,
	financialTransactionHandler *handler.FinancialTransactionHandler,
) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())

	// API路由组
	api := router.Group("/api")
	{
		// 客户路由
		customers := api.Group("/customers")
		{
			customers.GET("", customerHandler.List)
			customers.GET("/:id", customerHandler.Get)
			customers.POST("", customerHandler.Create)
			customers.PUT("/:id", customerHandler.Update)
			customers.DELETE("/:id", customerHandler.Delete)
		}

		// SKU分类路由
		skuCategories := api.Group("/sku-categories")
		{
			skuCategories.GET("", skuCategoryHandler.List)
			skuCategories.GET("/:id", skuCategoryHandler.Get)
			skuCategories.POST("", skuCategoryHandler.Create)
			skuCategories.PUT("/:id", skuCategoryHandler.Update)
			skuCategories.DELETE("/:id", skuCategoryHandler.Delete)
		}

		// SKU路由
		skus := api.Group("/skus")
		{
			skus.GET("", skuHandler.List)
			skus.GET("/all", skuHandler.ListAll)
			skus.GET("/:id", skuHandler.Get)
			skus.POST("", skuHandler.Create)
			skus.PUT("/:id", skuHandler.Update)
			skus.DELETE("/:id", skuHandler.Delete)
		}

		// 订单路由
		orders := api.Group("/orders")
		{
			orders.GET("", orderHandler.List)
			orders.GET("/:id", orderHandler.Get)
			orders.POST("", orderHandler.Create)
			orders.PUT("/:id", orderHandler.Update)
			orders.PUT("/:id/status", orderHandler.UpdateStatus)
			orders.DELETE("/:id", orderHandler.Delete)
			orders.POST("/:id/settle", orderHandler.Settle)
		}

		// 财务流水路由
		transactions := api.Group("/financial-transactions")
		{
			transactions.GET("", financialTransactionHandler.List)
			transactions.GET("/balance", financialTransactionHandler.GetBalance)
			transactions.POST("", financialTransactionHandler.Create)
		}
	}

	return router
}
