package service

import (
	"super-order-web/internal/model"

	"gorm.io/gorm"
)

// SKUService SKU服务
type SKUService struct {
	db *gorm.DB
}

// NewSKUService 创建SKU服务
func NewSKUService(db *gorm.DB) *SKUService {
	return &SKUService{db: db}
}

// List 获取SKU列表
func (s *SKUService) List(page, pageSize int) ([]model.SKU, int64, error) {
	var skus []model.SKU
	var total int64

	query := s.db.Where("is_deleted = ?", false)

	if err := query.Model(&model.SKU{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Category").Offset(offset).Limit(pageSize).Find(&skus).Error

	return skus, total, err
}

// GetByID 根据ID获取SKU
func (s *SKUService) GetByID(id int64) (*model.SKU, error) {
	var sku model.SKU
	err := s.db.Preload("Category").Where("id = ?", id).First(&sku).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
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

// Create 创建SKU
func (s *SKUService) Create(sku *model.SKU) error {
	return s.db.Create(sku).Error
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
func (s *SKUService) ListAll() ([]model.SKU, error) {
	var skus []model.SKU
	err := s.db.Where("is_deleted = ?", false).Find(&skus).Error
	return skus, err
}
