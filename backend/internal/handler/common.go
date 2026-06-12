package handler

import (
	"super-order-web/internal/service"
	"super-order-web/pkg/response"

	"github.com/gin-gonic/gin"
)

// CommonHandler 公共处理器
type CommonHandler struct {
	skuService *service.SKUService
}

// NewCommonHandler 创建公共处理器
func NewCommonHandler(skuService *service.SKUService) *CommonHandler {
	return &CommonHandler{skuService: skuService}
}

// GetImageResponse SKU图片响应
type GetImageResponse struct {
	SKUCode  string `json:"sku_code"`
	ImageURL string `json:"image_url,omitempty"`
	HasImage bool   `json:"has_image"`
}

// GetImage 获取SKU图片
// @Summary 获取SKU图片
// @Tags Common
// @Accept json
// @Produce json
// @Param skuCode path string true "SKU编码"
// @Success 200 {object} response.Response
// @Router /api/common/image/{skuCode} [get]
func (h *CommonHandler) GetImage(c *gin.Context) {
	skuCode := c.Param("skuCode")
	_, err := h.skuService.GetByCode(skuCode)
	if err != nil {
		response.Fail(c, "SKU不存在")
		return
	}

	// TODO: 从OSS获取图片URL
	// 目前先返回空，后续可以集成OSS
	resp := GetImageResponse{
		SKUCode:  skuCode,
		ImageURL: "",
		HasImage: false,
	}

	response.Success(c, resp)
}
