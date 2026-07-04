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
	service                      *service.OrderService
	financialTransactionService *service.FinancialTransactionService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(svc *service.OrderService, financialSvc *service.FinancialTransactionService) *OrderHandler {
	return &OrderHandler{
		service:                      svc,
		financialTransactionService: financialSvc,
	}
}

// ListRequest 获取订单列表请求
type ListOrderRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Status   string `form:"status"`
}

// OrderItemResponse 订单明细响应（扁平化SKU字段）
type OrderItemResponse struct {
	ID               int64   `json:"id"`
	OrderID          int64   `json:"order_id"`
	SKUID            int64   `json:"sku_id"`
	SKUCode          string  `json:"sku_code"`
	ProductName      string  `json:"product_name"`
	Spec             string  `json:"spec"`
	Unit             string  `json:"unit"`
	BoxSpec          string  `json:"box_spec"`
	BoxQuantity      int     `json:"box_quantity"`
	Quantity         int     `json:"quantity"`
	CostPrice        float64 `json:"cost_price"`
	SalePrice        float64 `json:"sale_price"`
	TotalCostAmount  float64 `json:"total_cost_amount"`
	TotalSaleAmount  float64 `json:"total_sale_amount"`
	SettledAmount    float64 `json:"settled_amount"`
}

// ToOrderItemResponse 将 model.OrderItem 转换为 OrderItemResponse
func ToOrderItemResponse(item model.OrderItem) OrderItemResponse {
	resp := OrderItemResponse{
		ID:               item.ID,
		OrderID:          item.OrderID,
		SKUID:            item.SKUID,
		SKUCode:          item.SKUCode,
		ProductName:      item.ProductName,
		Spec:             "",
		Unit:             "个",
		BoxSpec:          "",
		BoxQuantity:      1,
		Quantity:         item.Quantity,
		CostPrice:        item.CostPrice,
		SalePrice:        item.SalePrice,
		TotalCostAmount:  item.TotalCostAmount,
		TotalSaleAmount:  item.TotalSaleAmount,
		SettledAmount:    item.SettledAmount,
	}
	if item.SKU != nil {
		resp.Spec = item.SKU.Spec
		resp.Unit = item.SKU.Unit
		resp.BoxSpec = item.SKU.BoxSpec
		resp.BoxQuantity = item.SKU.BoxQuantity
	}
	return resp
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
// @Router /api/order/list [get]
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
// @Router /api/order/{id} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "订单不存在")
		return
	}
	response.Success(c, order)
}

// CreateRequest 创建订单请求
type CreateOrderRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	OrderDate  string `json:"order_date" binding:"required"`
	Remarks    string `json:"remarks"`
}

// Create 创建订单
// @Summary 创建订单
// @Tags Order
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "订单信息"
// @Success 200 {object} response.Response
// @Router /api/order [post]
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

	if err := h.service.Create(order); err != nil {
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
	IsSettled       bool    `json:"is_settled"`
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
// @Router /api/order/{id} [put]
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
		IsSettled:       req.IsSettled,
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

// UpdateStatus 更新订单状态（未在路由中使用）
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
// @Router /api/order/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// Settle 订单结算（未在路由中使用）
func (h *OrderHandler) Settle(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Settle(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// GetProcessingOrders 获取进行中的订单
// @Summary 获取进行中的订单
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/order/processing [get]
func (h *OrderHandler) GetProcessingOrders(c *gin.Context) {
	orders, err := h.service.GetByStatus([]string{"processing"})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, orders)
}

// GetUnsettledOrders 获取未结算的订单
// @Summary 获取未结算的订单
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/order/unsettled [get]
func (h *OrderHandler) GetUnsettledOrders(c *gin.Context) {
	orders, err := h.service.GetUnsettled()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, orders)
}

// GetItems 获取订单明细
// @Summary 获取订单明细
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response
// @Router /api/order/{id}/items [get]
func (h *OrderHandler) GetItems(c *gin.Context) {
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	items, err := h.service.GetItems(orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 获取最新余额
	lastBalance, err := h.financialTransactionService.GetLatestBalance()
	if err != nil {
		lastBalance = 0
	}

	// 转换为扁平化的响应结构
	resp := make([]OrderItemResponse, len(items))
	for i, item := range items {
		resp[i] = ToOrderItemResponse(item)
	}

	response.Success(c, map[string]interface{}{
		"items":       resp,
		"last_balance": lastBalance,
	})
}
