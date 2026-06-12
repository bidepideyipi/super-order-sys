package handler

import (
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/response"

	"github.com/gin-gonic/gin"
)

// SKUCategoryHandler SKU分类处理器
type SKUCategoryHandler struct {
	service *service.SKUCategoryService
}

// NewSKUCategoryHandler 创建SKU分类处理器
func NewSKUCategoryHandler(svc *service.SKUCategoryService) *SKUCategoryHandler {
	return &SKUCategoryHandler{service: svc}
}

// List 获取分类列表
// @Summary 获取SKU分类列表
// @Tags SKUCategory
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/category/list [get]
func (h *SKUCategoryHandler) List(c *gin.Context) {
	categories, err := h.service.List()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, categories)
}

// Get 获取分类详情
// @Summary 获取SKU分类详情
// @Tags SKUCategory
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} response.Response
// @Router /api/category/{id} [get]
func (h *SKUCategoryHandler) Get(c *gin.Context) {
	id := c.Param("id")
	category, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "分类不存在")
		return
	}
	response.Success(c, category)
}

// CreateRequest 创建分类请求
type CreateCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"`
}

// Create 创建分类
// @Summary 创建SKU分类
// @Tags SKUCategory
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/category [post]
func (h *SKUCategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	category := &model.SKUCategory{
		CategoryName: req.CategoryName,
	}

	if err := h.service.Create(category); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, category)
}

// UpdateRequest 更新分类请求
type UpdateCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"`
}

// Update 更新分类
// @Summary 更新SKU分类
// @Tags SKUCategory
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Param request body UpdateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/category/{id} [put]
func (h *SKUCategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	category := &model.SKUCategory{
		CategoryID:   id,
		CategoryName: req.CategoryName,
	}

	if err := h.service.Update(category); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// Delete 删除分类
// @Summary 删除SKU分类
// @Tags SKUCategory
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} response.Response
// @Router /api/category/{id} [delete]
func (h *SKUCategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}
