package handler

import (
	"strconv"
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/response"
	"super-order-web/pkg/util"

	"github.com/gin-gonic/gin"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

// ListRequest 获取订单列表请求
type ListOrderRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Status   string `form:"status"`
}

// List 获取订单列表
// @Summary 获取订单列表
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "订单状态"
// @Success 200 {object} response.Response
// @Router /api/orders [get]
func (h *OrderHandler) List(c *gin.Context) {
	var req ListOrderRequest
	req.Page = 1
	req.PageSize = 10
	req.Status = "all"
	c.ShouldBindQuery(&req)

	orders, total, err := h.service.List(req.Page, req.PageSize, req.Status)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, util.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Data:     orders,
	})
}

// Get 获取订单详情
// @Summary 获取订单详情
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response
// @Router /api/orders/{id} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "订单不存在")
		return
	}
	response.Success(c, order)
}

// OrderItemRequest 订单明细请求
type OrderItemRequest struct {
	SKUID          int64   `json:"sku_id"`
	SKUCode        string  `json:"sku_code" binding:"required"`
	ProductName    string  `json:"product_name" binding:"required"`
	Quantity       int     `json:"quantity" binding:"required,min=1"`
	CostPrice      float64 `json:"cost_price" binding:"required"`
	SalePrice      float64 `json:"sale_price" binding:"required"`
	TotalCostAmount float64 `json:"total_cost_amount" binding:"required"`
	TotalSaleAmount float64 `json:"total_sale_amount" binding:"required"`
}

// CreateRequest 创建订单请求
type CreateOrderRequest struct {
	CustomerID string           `json:"customer_id" binding:"required"`
	OrderDate  string           `json:"order_date" binding:"required"`
	Remarks    string           `json:"remarks"`
	Items      []OrderItemRequest `json:"items" binding:"required,min=1"`
}

// Create 创建订单
// @Summary 创建订单
// @Tags Order
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "订单信息"
// @Success 200 {object} response.Response
// @Router /api/orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	order := &model.Order{
		CustomerID: req.CustomerID,
		OrderDate:  req.OrderDate,
		Remarks:    req.Remarks,
		Status:     string(model.OrderStatusPending),
	}

	items := make([]model.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = model.OrderItem{
			SKUID:           item.SKUID,
			SKUCode:         item.SKUCode,
			ProductName:     item.ProductName,
			Quantity:        item.Quantity,
			CostPrice:       item.CostPrice,
			SalePrice:       item.SalePrice,
			TotalCostAmount: item.TotalCostAmount,
			TotalSaleAmount: item.TotalSaleAmount,
		}
	}

	if err := h.service.Create(order, items); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, order)
}

// UpdateRequest 更新订单请求
type UpdateOrderRequest struct {
	CustomerID      string  `json:"customer_id" binding:"required"`
	OrderDate       string  `json:"order_date" binding:"required"`
	Status          string  `json:"status"`
	TotalCostAmount float64 `json:"total_cost_amount"`
	TotalSaleAmount float64 `json:"total_sale_amount"`
	Remarks         string  `json:"remarks"`
}

// Update 更新订单
// @Summary 更新订单
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body UpdateOrderRequest true "订单信息"
// @Success 200 {object} response.Response
// @Router /api/orders/{id} [put]
func (h *OrderHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	order := &model.Order{
		ID:              id,
		CustomerID:      req.CustomerID,
		OrderDate:       req.OrderDate,
		Status:          req.Status,
		TotalCostAmount: req.TotalCostAmount,
		TotalSaleAmount: req.TotalSaleAmount,
		Remarks:         req.Remarks,
	}

	if err := h.service.Update(order); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdateStatusRequest 更新订单状态请求
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus 更新订单状态
// @Summary 更新订单状态
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body UpdateStatusRequest true "状态信息"
// @Success 200 {object} response.Response
// @Router /api/orders/{id}/status [put]
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// Delete 删除订单
// @Summary 删除订单
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response
// @Router /api/orders/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// Settle 订单结算
// @Summary 订单结算
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response
// @Router /api/orders/{id}/settle [post]
func (h *OrderHandler) Settle(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Settle(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}
