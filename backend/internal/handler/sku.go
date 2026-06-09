package handler

import (
	"strconv"
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/response"
	"super-order-web/pkg/util"

	"github.com/gin-gonic/gin"
)

// SKUHandler SKU处理器
type SKUHandler struct {
	service *service.SKUService
}

// NewSKUHandler 创建SKU处理器
func NewSKUHandler(svc *service.SKUService) *SKUHandler {
	return &SKUHandler{service: svc}
}

// ListRequest 获取SKU列表请求
type ListSKURequest struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

// List 获取SKU列表
// @Summary 获取SKU列表
// @Tags SKU
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/skus [get]
func (h *SKUHandler) List(c *gin.Context) {
	var req ListSKURequest
	req.Page = 1
	req.PageSize = 10
	c.ShouldBindQuery(&req)

	skus, total, err := h.service.List(req.Page, req.PageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, util.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Data:     skus,
	})
}

// ListAll 获取所有SKU
// @Summary 获取所有SKU（不分页）
// @Tags SKU
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/skus/all [get]
func (h *SKUHandler) ListAll(c *gin.Context) {
	skus, err := h.service.ListAll()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, skus)
}

// Get 获取SKU详情
// @Summary 获取SKU详情
// @Tags SKU
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} response.Response
// @Router /api/skus/{id} [get]
func (h *SKUHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sku, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "SKU不存在")
		return
	}
	response.Success(c, sku)
}

// CreateRequest 创建SKU请求
type CreateSKURequest struct {
	SKUCode     string  `json:"sku_code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Spec        string  `json:"spec"`
	Unit        string  `json:"unit"`
	CategoryID  string  `json:"category_id" binding:"required"`
	BoxSpec     string  `json:"box_spec"`
	BoxQuantity int     `json:"box_quantity"`
	CostPrice   float64 `json:"cost_price"`
	SalePrice   float64 `json:"sale_price"`
}

// Create 创建SKU
// @Summary 创建SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param request body CreateSKURequest true "SKU信息"
// @Success 200 {object} response.Response
// @Router /api/skus [post]
func (h *SKUHandler) Create(c *gin.Context) {
	var req CreateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	sku := &model.SKU{
		SKUCode:     req.SKUCode,
		Name:        req.Name,
		Description: req.Description,
		Spec:        req.Spec,
		Unit:        req.Unit,
		CategoryID:  req.CategoryID,
		BoxSpec:     req.BoxSpec,
		BoxQuantity: req.BoxQuantity,
		CostPrice:   req.CostPrice,
		SalePrice:   req.SalePrice,
	}

	if sku.Unit == "" {
		sku.Unit = "个"
	}
	if sku.BoxQuantity == 0 {
		sku.BoxQuantity = 1
	}

	if err := h.service.Create(sku); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, sku)
}

// UpdateRequest 更新SKU请求
type UpdateSKURequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Spec        string  `json:"spec"`
	Unit        string  `json:"unit"`
	CategoryID  string  `json:"category_id" binding:"required"`
	BoxSpec     string  `json:"box_spec"`
	BoxQuantity int     `json:"box_quantity"`
	CostPrice   float64 `json:"cost_price"`
	SalePrice   float64 `json:"sale_price"`
}

// Update 更新SKU
// @Summary 更新SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Param request body UpdateSKURequest true "SKU信息"
// @Success 200 {object} response.Response
// @Router /api/skus/{id} [put]
func (h *SKUHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	sku := &model.SKU{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Spec:        req.Spec,
		Unit:        req.Unit,
		CategoryID:  req.CategoryID,
		BoxSpec:     req.BoxSpec,
		BoxQuantity: req.BoxQuantity,
		CostPrice:   req.CostPrice,
		SalePrice:   req.SalePrice,
	}

	if err := h.service.Update(sku); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// Delete 删除SKU
// @Summary 删除SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} response.Response
// @Router /api/skus/{id} [delete]
func (h *SKUHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}
