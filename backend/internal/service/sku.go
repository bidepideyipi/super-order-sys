package service

import (
	"fmt"
	"super-order-web/internal/model"
	"super-order-web/pkg/oss"

	"gorm.io/gorm"
)

// SKUWithImage SKU带图片字段的响应结构
type SKUWithImage struct {
	model.SKU
	CategoryName string `json:"category_name"`
	Image        string `json:"image"`
}

// toSKUWithImage 将SKU转换为带图片字段的结构
func toSKUWithImage(sku model.SKU) SKUWithImage {
	categoryName := ""
	if sku.Category != nil {
		categoryName = sku.Category.CategoryName
	}
	return SKUWithImage{
		SKU:          sku,
		CategoryName: categoryName,
		Image:        oss.GetImageURL(sku.SKUCode),
	}
}

// toSKUWithImageSlice 批量转换
func toSKUWithImageSlice(skus []model.SKU) []SKUWithImage {
	result := make([]SKUWithImage, len(skus))
	for i, sku := range skus {
		result[i] = toSKUWithImage(sku)
	}
	return result
}

// SKUService SKU服务
type SKUService struct {
	db *gorm.DB
}

// NewSKUService 创建SKU服务
func NewSKUService(db *gorm.DB) *SKUService {
	return &SKUService{db: db}
}

// List 获取SKU列表
func (s *SKUService) List(page, pageSize int) ([]SKUWithImage, int64, error) {
	var skus []model.SKU
	var total int64

	query := s.db.Where("is_deleted = ?", false)

	if err := query.Model(&model.SKU{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Category").Offset(offset).Limit(pageSize).Find(&skus).Error
	if err != nil {
		return nil, 0, err
	}

	return toSKUWithImageSlice(skus), total, nil
}

// GetByID 根据ID获取SKU
func (s *SKUService) GetByID(id int64) (*SKUWithImage, error) {
	var sku model.SKU
	err := s.db.Preload("Category").Where("id = ?", id).First(&sku).Error
	if err != nil {
		return nil, err
	}
	result := toSKUWithImage(sku)
	return &result, nil
}

// GetByCode 根据编码获取SKU
func (s *SKUService) GetByCode(code string) (*model.SKU, error) {
	var sku model.SKU
	err := s.db.Preload("Category").Where("sku_code = ?", code).First(&sku).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

// GenerateSKUCode 生成SKU编码
// 规则：分类ID(2位) + 序号(4位，不足左补0)
func (s *SKUService) GenerateSKUCode(categoryID string) (string, error) {
	// 获取该分类下最大的ID
	var maxID int64
	err := s.db.Model(&model.SKU{}).
		Where("category_id = ?", categoryID).
		Select("COALESCE(MAX(id), 0)").
		Scan(&maxID).Error
	if err != nil {
		return "", err
	}

	// 序号 = 最大ID + 1
	seq := maxID + 1

	// 分类ID取前2位作为分类编码
	categoryCode := categoryID
	if len(categoryCode) > 2 {
		categoryCode = categoryCode[:2]
	}
	if len(categoryCode) < 2 {
		categoryCode = "0" + categoryCode
	}

	// 序号补0到4位
	seqStr := fmt.Sprintf("%04d", seq)

	skuCode := categoryCode + seqStr
	return skuCode, nil
}

// Create 创建SKU
func (s *SKUService) Create(sku *model.SKU) error {
	// 如果SKU编码为空，自动生成
	if sku.SKUCode == "" {
		skuCode, err := s.GenerateSKUCode(sku.CategoryID)
		if err != nil {
			return err
		}
		sku.SKUCode = skuCode
	}

	err := s.db.Create(sku).Error
	if err != nil {
		fmt.Printf("Create SKU error: %v, SKU: %+v\n", err, sku)
	} else {
		fmt.Printf("Create SKU success: ID=%d, SKUCode=%s\n", sku.ID, sku.SKUCode)
	}
	return err
}

// Update 更新SKU
func (s *SKUService) Update(sku *model.SKU) error {
	return s.db.Model(&model.SKU{}).
		Where("id = ?", sku.ID).
		Updates(map[string]interface{}{
			"name":          sku.Name,
			"description":   sku.Description,
			"spec":          sku.Spec,
			"unit":          sku.Unit,
			"category_id":   sku.CategoryID,
			"box_spec":      sku.BoxSpec,
			"box_quantity":  sku.BoxQuantity,
			"cost_price":    sku.CostPrice,
			"sale_price":    sku.SalePrice,
		}).Error
}

// Delete 软删除SKU
func (s *SKUService) Delete(id int64) error {
	return s.db.Model(&model.SKU{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// ListAll 获取所有SKU（不分页）
func (s *SKUService) ListAll() ([]SKUWithImage, error) {
	var skus []model.SKU
	err := s.db.Where("is_deleted = ?", false).Preload("Category").Find(&skus).Error
	if err != nil {
		return nil, err
	}
	return toSKUWithImageSlice(skus), nil
}

// Search 搜索SKU
func (s *SKUService) Search(keyword string) ([]SKUWithImage, error) {
	var skus []model.SKU
	query := s.db.Where("is_deleted = ?", false)
	if keyword != "" {
		query = query.Where("sku_code LIKE ? OR name LIKE ? OR spec LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := query.Preload("Category").Find(&skus).Error
	if err != nil {
		return nil, err
	}
	return toSKUWithImageSlice(skus), nil
}

// SearchPaginated 分页搜索SKU
func (s *SKUService) SearchPaginated(keyword string, page, pageSize int) ([]SKUWithImage, int64, error) {
	var skus []model.SKU
	var total int64

	query := s.db.Where("is_deleted = ?", false)
	if keyword != "" {
		query = query.Where("sku_code LIKE ? OR name LIKE ? OR spec LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Model(&model.SKU{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Category").Offset(offset).Limit(pageSize).Find(&skus).Error
	if err != nil {
		return nil, 0, err
	}

	return toSKUWithImageSlice(skus), total, nil
}

// SearchWithCategory 按分类搜索SKU
func (s *SKUService) SearchWithCategory(keyword, categoryID string) ([]SKUWithImage, error) {
	var skus []model.SKU
	query := s.db.Where("is_deleted = ?", false)

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("sku_code LIKE ? OR name LIKE ? OR spec LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Preload("Category").Find(&skus).Error
	if err != nil {
		return nil, err
	}
	return toSKUWithImageSlice(skus), nil
}
