package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/wmsflow/wmsflow/internal/config"
	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
)

// Connect 连接 PostgreSQL 并自动迁移表结构。
func Connect(cfg *config.Config, slogger *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	gormLogger := logger.New(
		&slogWriter{slogger},
		logger.Config{SlowThreshold: 500 * time.Millisecond, LogLevel: logger.Warn},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	slogger.Info(constants.LogDBConnected, "host", cfg.DBHost, "db", cfg.DBName)
	if err := db.AutoMigrate(
		&model.User{},
		&model.Owner{},
		&model.Product{},
		&model.BinLocation{},
		&model.Inventory{},
		&model.InboundOrder{},
		&model.InboundItem{},
		&model.OutboundOrder{},
		&model.OutboundItem{},
		&model.OperationLog{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	slogger.Info(constants.LogMigrateCompleted, "tables", 10)
	return db, nil
}

// slogWriter 将 GORM 日志桥接到 slog。
type slogWriter struct {
	logger *slog.Logger
}

func (w *slogWriter) Printf(format string, args ...any) {
	w.logger.Debug(fmt.Sprintf(format, args...))
}
