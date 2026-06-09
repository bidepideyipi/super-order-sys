package handler

import (
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/response"

	"github.com/gin-gonic/gin"
)

// CustomerHandler 客户处理器
type CustomerHandler struct {
	service *service.CustomerService
}

// NewCustomerHandler 创建客户处理器
func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: svc}
}

// List 获取客户列表
// @Summary 获取客户列表
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/customers [get]
func (h *CustomerHandler) List(c *gin.Context) {
	customers, err := h.service.List()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, customers)
}

// Get 获取客户详情
// @Summary 获取客户详情
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "客户ID"
// @Success 200 {object} response.Response
// @Router /api/customers/{id} [get]
func (h *CustomerHandler) Get(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "客户不存在")
		return
	}
	response.Success(c, customer)
}

// CreateRequest 创建客户请求
type CreateRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
}

// Create 创建客户
// @Summary 创建客户
// @Tags Customer
// @Accept json
// @Produce json
// @Param request body CreateRequest true "客户信息"
// @Success 200 {object} response.Response
// @Router /api/customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	customer := &model.Customer{
		CustomerName: req.CustomerName,
	}

	if err := h.service.Create(customer); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, customer)
}

// UpdateRequest 更新客户请求
type UpdateRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
}

// Update 更新客户
// @Summary 更新客户
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "客户ID"
// @Param request body UpdateRequest true "客户信息"
// @Success 200 {object} response.Response
// @Router /api/customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	customer := &model.Customer{
		CustomerID:   id,
		CustomerName: req.CustomerName,
	}

	if err := h.service.Update(customer); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// Delete 删除客户
// @Summary 删除客户
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "客户ID"
// @Success 200 {object} response.Response
// @Router /api/customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}
