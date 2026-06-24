package database

import (
	"fmt"
	"super-order-web/internal/config"

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

	// 获取底层的sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置SQLite参数
	sqlDB.Exec("PRAGMA foreign_keys = ON")
	sqlDB.Exec("PRAGMA journal_mode = WAL")
	sqlDB.Exec("PRAGMA synchronous = NORMAL")

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
