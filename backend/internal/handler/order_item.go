package handler

import (
	"strconv"
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/response"

	"github.com/gin-gonic/gin"
)

// OrderItemHandler 订单明细处理器
type OrderItemHandler struct {
	service *service.OrderItemService
}

// NewOrderItemHandler 创建订单明细处理器
func NewOrderItemHandler(svc *service.OrderItemService) *OrderItemHandler {
	return &OrderItemHandler{service: svc}
}

// Get 获取订单明细
// @Summary 获取订单明细
// @Tags OrderItem
// @Accept json
// @Produce json
// @Param id path int true "订单明细ID"
// @Success 200 {object} response.Response
// @Router /api/order-item/{id} [get]
func (h *OrderItemHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "订单明细不存在")
		return
	}
	response.Success(c, item)
}

// CreateRequest 创建订单明细请求
type CreateOrderItemRequest struct {
	OrderID         int64   `json:"order_id" binding:"required"`
	SKUID           int64   `json:"sku_id"`
	SKUCode         string  `json:"sku_code" binding:"required"`
	ProductName     string  `json:"product_name" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required,min=1"`
	CostPrice       float64 `json:"cost_price" binding:"required"`
	SalePrice       float64 `json:"sale_price" binding:"required"`
	TotalCostAmount float64 `json:"total_cost_amount" binding:"required"`
	TotalSaleAmount float64 `json:"total_sale_amount" binding:"required"`
}

// Create 创建订单明细
// @Summary 创建订单明细
// @Tags OrderItem
// @Accept json
// @Produce json
// @Param request body CreateOrderItemRequest true "订单明细信息"
// @Success 200 {object} response.Response
// @Router /api/order-item [post]
func (h *OrderItemHandler) Create(c *gin.Context) {
	var req CreateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	item := &model.OrderItem{
		OrderID:         req.OrderID,
		SKUID:           req.SKUID,
		SKUCode:         req.SKUCode,
		ProductName:     req.ProductName,
		Quantity:        req.Quantity,
		CostPrice:       req.CostPrice,
		SalePrice:       req.SalePrice,
		TotalCostAmount: req.TotalCostAmount,
		TotalSaleAmount: req.TotalSaleAmount,
	}

	if err := h.service.Create(item); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, item)
}

// UpdateRequest 更新订单明细请求
type UpdateOrderItemRequest struct {
	SKUID           int64   `json:"sku_id"`
	SKUCode         string  `json:"sku_code" binding:"required"`
	ProductName     string  `json:"product_name" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required,min=1"`
	CostPrice       float64 `json:"cost_price" binding:"required"`
	SalePrice       float64 `json:"sale_price" binding:"required"`
	TotalCostAmount float64 `json:"total_cost_amount" binding:"required"`
	TotalSaleAmount float64 `json:"total_sale_amount" binding:"required"`
	SettledAmount   float64 `json:"settled_amount"`
}

// Update 更新订单明细
// @Summary 更新订单明细
// @Tags OrderItem
// @Accept json
// @Produce json
// @Param id path int true "订单明细ID"
// @Param request body UpdateOrderItemRequest true "订单明细信息"
// @Success 200 {object} response.Response
// @Router /api/order-item/{id} [put]
func (h *OrderItemHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	item := &model.OrderItem{
		ID:              id,
		SKUID:           req.SKUID,
		SKUCode:         req.SKUCode,
		ProductName:     req.ProductName,
		Quantity:        req.Quantity,
		CostPrice:       req.CostPrice,
		SalePrice:       req.SalePrice,
		TotalCostAmount: req.TotalCostAmount,
		TotalSaleAmount: req.TotalSaleAmount,
		SettledAmount:   req.SettledAmount,
	}

	if err := h.service.Update(item); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// Delete 删除订单明细
// @Summary 删除订单明细
// @Tags OrderItem
// @Accept json
// @Produce json
// @Param id path int true "订单明细ID"
// @Success 200 {object} response.Response
// @Router /api/order-item/{id} [delete]
func (h *OrderItemHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}
