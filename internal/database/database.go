package database

import (
	"fmt"
	"log"
	"time"

	"crud-api/internal/config"
	"crud-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	var db *gorm.DB
	var err error

	// Thử kết nối lại nếu DB chưa sẵn sàng (đặc biệt khi chạy qua Docker)
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		log.Printf("[Lần %d/%d] Đang chờ kết nối tới Database... Lỗi: %v", i, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới cơ sở dữ liệu sau %d lần thử: %w", maxRetries, err)
	}

	// Cấu hình Connection Pool cho PostgreSQL
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("lỗi cấu hình connection pool: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Kết nối PostgreSQL thành công!")

	// Tự động migrate schema
	if err := db.AutoMigrate(
		&model.Product{},
		&model.User{},
		&model.Story{},
		&model.SystemCategory{},
		&model.SystemCategoryItem{},
		&model.Location{},
		&model.Faction{},
		&model.Character{},
		&model.Plot{},
		&model.Item{},
	); err != nil {
		return nil, fmt.Errorf("lỗi khi AutoMigrate các bảng dữ liệu: %w", err)
	}
	log.Println("Database AutoMigrate hoàn tất!")

	return db, nil
}
