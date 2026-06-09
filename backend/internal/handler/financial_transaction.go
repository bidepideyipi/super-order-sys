package handler

import (
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/response"
	"super-order-web/pkg/util"

	"github.com/gin-gonic/gin"
)

// FinancialTransactionHandler 财务流水处理器
type FinancialTransactionHandler struct {
	service *service.FinancialTransactionService
}

// NewFinancialTransactionHandler 创建财务流水处理器
func NewFinancialTransactionHandler(svc *service.FinancialTransactionService) *FinancialTransactionHandler {
	return &FinancialTransactionHandler{service: svc}
}

// ListRequest 获取财务流水列表请求
type ListTransactionRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Category string `form:"category"`
}

// List 获取财务流水列表
// @Summary 获取财务流水列表
// @Tags FinancialTransaction
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param category query string false "分类"
// @Success 200 {object} response.Response
// @Router /api/financial-transactions [get]
func (h *FinancialTransactionHandler) List(c *gin.Context) {
	var req ListTransactionRequest
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

// CreateRequest 创建财务流水请求
type CreateTransactionRequest struct {
	Category     string  `json:"category" binding:"required"`
	Description  string  `json:"description"`
	AmountChange float64 `json:"amount_change" binding:"required"`
}

// Create 创建财务流水
// @Summary 创建财务流水
// @Tags FinancialTransaction
// @Accept json
// @Produce json
// @Param request body CreateTransactionRequest true "财务流水信息"
// @Success 200 {object} response.Response
// @Router /api/financial-transactions [post]
func (h *FinancialTransactionHandler) Create(c *gin.Context) {
	var req CreateTransactionRequest
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

// GetBalance 获取当前余额
// @Summary 获取当前余额
// @Tags FinancialTransaction
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/financial-transactions/balance [get]
func (h *FinancialTransactionHandler) GetBalance(c *gin.Context) {
	balance, err := h.service.GetLatestBalance()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, map[string]float64{"balance": balance})
}
