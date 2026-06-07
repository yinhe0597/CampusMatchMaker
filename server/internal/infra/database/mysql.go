package database

import (
	"campus_collab/internal/infra/config"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// InitMySQL 初始化 MySQL 连接
func InitMySQL(cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.DSN()
	log.Info("正在连接 MySQL", zap.String("host", cfg.Host), zap.Int("port", cfg.Port), zap.String("database", cfg.Name))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("MySQL 连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 验证连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL Ping 失败: %w", err)
	}

	log.Info("MySQL 连接成功")
	return db, nil
}
