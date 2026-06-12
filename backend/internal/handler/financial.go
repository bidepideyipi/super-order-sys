package handler

import (
	"strconv"
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/response"
	"super-order-web/pkg/util"

	"github.com/gin-gonic/gin"
)

// FinancialHandler 财务处理器
type FinancialHandler struct {
	service *service.FinancialTransactionService
}

// NewFinancialHandler 创建财务处理器
func NewFinancialHandler(svc *service.FinancialTransactionService) *FinancialHandler {
	return &FinancialHandler{service: svc}
}

// ListRequest 获取财务列表请求
type ListFinancialRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Category string `form:"category"`
}

// List 获取财务列表
// @Summary 获取财务列表
// @Tags Financial
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param category query string false "分类"
// @Success 200 {object} response.Response
// @Router /api/financial/list [get]
func (h *FinancialHandler) List(c *gin.Context) {
	var req ListFinancialRequest
	req.Page = 1
	req.PageSize = 10
	req.Category = "all"
	c.ShouldBindQuery(&req)

	transactions, total, err := h.service.List(req.Page, req.PageSize, req.Category)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, util.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Data:     transactions,
	})
}

// Get 获取财务详情
// @Summary 获取财务详情
// @Tags Financial
// @Accept json
// @Produce json
// @Param id path int true "财务ID"
// @Success 200 {object} response.Response
// @Router /api/financial/{id} [get]
func (h *FinancialHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	transaction, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "记录不存在")
		return
	}
	response.Success(c, transaction)
}

// CreateRequest 创建财务请求
type CreateFinancialRequest struct {
	Category     string  `json:"category" binding:"required"`
	Description  string  `json:"description"`
	AmountChange float64 `json:"amount_change" binding:"required"`
}

// Create 创建财务
// @Summary 创建财务
// @Tags Financial
// @Accept json
// @Produce json
// @Param request body CreateFinancialRequest true "财务信息"
// @Success 200 {object} response.Response
// @Router /api/financial [post]
func (h *FinancialHandler) Create(c *gin.Context) {
	var req CreateFinancialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 获取最新余额
	latestBalance, err := h.service.GetLatestBalance()
	if err != nil {
		response.Error(c, err)
		return
	}

	transaction := &model.FinancialTransaction{
		Category:     req.Category,
		Description:  req.Description,
		AmountChange: req.AmountChange,
		Balance:      latestBalance + req.AmountChange,
	}

	if err := h.service.Create(transaction); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, transaction)
}

// UpdateRequest 更新财务请求
type UpdateFinancialRequest struct {
	Category     string  `json:"category" binding:"required"`
	Description  string  `json:"description"`
	AmountChange float64 `json:"amount_change" binding:"required"`
	Balance      float64 `json:"balance" binding:"required"`
	IsSettled    bool    `json:"is_settled"`
}

// Update 更新财务
// @Summary 更新财务
// @Tags Financial
// @Accept json
// @Produce json
// @Param id path int true "财务ID"
// @Param request body UpdateFinancialRequest true "财务信息"
// @Success 200 {object} response.Response
// @Router /api/financial/{id} [put]
func (h *FinancialHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateFinancialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	transaction := &model.FinancialTransaction{
		ID:           id,
		Category:     req.Category,
		Description:  req.Description,
		AmountChange: req.AmountChange,
		Balance:      req.Balance,
		IsSettled:    req.IsSettled,
	}

	if err := h.service.Update(transaction); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// Delete 删除财务
// @Summary 删除财务
// @Tags Financial
// @Accept json
// @Produce json
// @Param id path int true "财务ID"
// @Success 200 {object} response.Response
// @Router /api/financial/{id} [delete]
func (h *FinancialHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// GetBalance 获取余额
// @Summary 获取余额
// @Tags Financial
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/financial/balance [get]
func (h *FinancialHandler) GetBalance(c *gin.Context) {
	balance, err := h.service.GetLatestBalance()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, map[string]float64{"balance": balance})
}
