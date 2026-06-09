package database

import (
	"fmt"
	"super-order-web/internal/config"
	"super-order-web/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize 初始化数据库连接
func Initialize(cfg *config.DatabaseConfig) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(cfg.SQLite), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.SKUCategory{},
		&model.Customer{},
		&model.FinancialTransaction{},
		&model.SKU{},
		&model.Order{},
		&model.OrderItem{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
