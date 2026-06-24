package handler

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"super-order-web/internal/model"
	"super-order-web/internal/service"
	"super-order-web/pkg/oss"
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

// List 获取SKU列表（分页，未在路由中使用）
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
// @Router /api/sku/list [get]
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
// @Router /api/sku/{id} [get]
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
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Spec        string  `json:"spec"`
	Unit        string  `json:"unit"`
	CategoryID  string  `json:"category_id" binding:"required"`
	BoxSpec     string  `json:"box_spec"`
	BoxQuantity int     `json:"box_quantity"`
	CostPrice   float64 `json:"cost_price"`
	SalePrice   float64 `json:"sale_price"`
	ImageBase64 string  `json:"image_base64"`
}

// Create 创建SKU
// @Summary 创建SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param request body CreateSKURequest true "SKU信息"
// @Success 200 {object} response.Response
// @Router /api/sku [post]
func (h *SKUHandler) Create(c *gin.Context) {
	var req CreateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	sku := &model.SKU{
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

	// 先生成SKU编码
	skuCode, err := h.service.GenerateSKUCode(req.CategoryID)
	if err != nil {
		response.Error(c, err)
		return
	}
	sku.SKUCode = skuCode

	// 上传图片到OSS
	imageURL := ""
	if req.ImageBase64 != "" {
		url, err := uploadImageToOSS(sku.SKUCode, req.ImageBase64)
		if err != nil {
			response.Fail(c, "图片上传失败: "+err.Error())
			return
		}
		imageURL = url
	}

	if err := h.service.Create(sku); err != nil {
		response.Error(c, err)
		return
	}

	// 返回包含图片URL的数据
	result := map[string]interface{}{
		"id":           sku.ID,
		"sku_code":     sku.SKUCode,
		"name":         sku.Name,
		"image":        imageURL,
		"category_name": "",
	}
	response.Success(c, result)
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
	ImageBase64 string  `json:"image_base64"`
}

// Update 更新SKU
// @Summary 更新SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Param request body UpdateSKURequest true "SKU信息"
// @Success 200 {object} response.Response
// @Router /api/sku/{id} [put]
func (h *SKUHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 获取SKU信息以便上传图片
	skuInfo, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, "SKU不存在")
		return
	}

	// 上传图片到OSS
	if req.ImageBase64 != "" {
		_, err := uploadImageToOSS(skuInfo.SKUCode, req.ImageBase64)
		if err != nil {
			response.Fail(c, "图片上传失败: "+err.Error())
			return
		}
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

// uploadImageToOSS 上传图片到OSS
func uploadImageToOSS(skuCode, imageBase64 string) (string, error) {
	// 解码base64
	base64Data := imageBase64
	if strings.HasPrefix(base64Data, "data:image/") {
		// 移除 data:image/xxx;base64, 前缀
		parts := strings.SplitN(base64Data, ",", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	// 解码base64
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// 上传到OSS，路径为 sku/{sku_code}.jpeg
	bucketName := oss.GetBucketName()
	if bucketName == "" {
		return "", nil // OSS未配置，返回空URL
	}

	objectName := fmt.Sprintf("sku/%s.jpeg", skuCode)
	url, err := oss.UploadBytes(bucketName, objectName, imageBytes)
	if err != nil {
		return "", err
	}

	return url, nil
}

// Delete 删除SKU
// @Summary 删除SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} response.Response
// @Router /api/sku/{id} [delete]
func (h *SKUHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// ListPaginated 分页获取SKU列表
// @Summary 分页获取SKU列表
// @Tags SKU
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/sku/list-paginated [get]
func (h *SKUHandler) ListPaginated(c *gin.Context) {
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

// Search 搜索SKU
// @Summary 搜索SKU（不分页）
// @Tags SKU
// @Accept json
// @Produce json
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response
// @Router /api/sku/search [get]
func (h *SKUHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	skus, err := h.service.Search(keyword)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, skus)
}

// SearchPaginated 分页搜索SKU
// @Summary 分页搜索SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param keyword query string true "关键词"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/sku/search-paginated [get]
func (h *SKUHandler) SearchPaginated(c *gin.Context) {
	keyword := c.Query("keyword")
	var req ListSKURequest
	req.Page = 1
	req.PageSize = 10
	c.ShouldBindQuery(&req)

	skus, total, err := h.service.SearchPaginated(keyword, req.Page, req.PageSize)
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

// SearchWithCategory 按分类搜索SKU
// @Summary 按分类搜索SKU
// @Tags SKU
// @Accept json
// @Produce json
// @Param keyword query string false "关键词"
// @Param category_id query string false "分类ID"
// @Success 200 {object} response.Response
// @Router /api/sku/search-with-category [get]
func (h *SKUHandler) SearchWithCategory(c *gin.Context) {
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")
	skus, err := h.service.SearchWithCategory(keyword, categoryID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, skus)
}
