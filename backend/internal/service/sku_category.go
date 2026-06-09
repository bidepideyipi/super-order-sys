package service

import (
	"super-order-web/internal/model"
	"super-order-web/pkg/util"

	"gorm.io/gorm"
)

// SKUCategoryService SKU分类服务
type SKUCategoryService struct {
	db *gorm.DB
}

// NewSKUCategoryService 创建SKU分类服务
func NewSKUCategoryService(db *gorm.DB) *SKUCategoryService {
	return &SKUCategoryService{db: db}
}

// List 获取分类列表
func (s *SKUCategoryService) List() ([]model.SKUCategory, error) {
	var categories []model.SKUCategory
	err := s.db.Find(&categories).Error
	return categories, err
}

// GetByID 根据ID获取分类
func (s *SKUCategoryService) GetByID(id string) (*model.SKUCategory, error) {
	var category model.SKUCategory
	err := s.db.Where("category_id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (s *SKUCategoryService) Create(category *model.SKUCategory) error {
	category.CategoryID = util.GenerateCategoryID()
	return s.db.Create(category).Error
}

// Update 更新分类
func (s *SKUCategoryService) Update(category *model.SKUCategory) error {
	return s.db.Model(&model.SKUCategory{}).
		Where("category_id = ?", category.CategoryID).
		Updates(map[string]interface{}{
			"category_name": category.CategoryName,
		}).Error
}

// Delete 删除分类
func (s *SKUCategoryService) Delete(id string) error {
	return s.db.Where("category_id = ?", id).Delete(&model.SKUCategory{}).Error
}
