package database

import (
	"errors"
	"log"

	"crud-api/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdmin tự động tạo tài khoản Admin mặc định nếu chưa tồn tại
func SeedAdmin(db *gorm.DB) error {
	adminEmail := "admin@example.com"
	adminUsername := "admin"
	adminPassword := "admin123"

	var existingUser model.User
	err := db.Where("email = ? OR role = ?", adminEmail, model.RoleAdmin).First(&existingUser).Error
	if err == nil {
		log.Printf("[Seeder] Tài khoản Admin đã tồn tại (Username: %s, Email: %s). Bỏ qua seed.", existingUser.Username, existingUser.Email)
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.User{
		Username: adminUsername,
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     model.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("[Seeder] Đã tạo thành công tài khoản Admin mặc định: Email='%s' | Mật khẩu='%s' | Role='%s'", adminEmail, adminPassword, model.RoleAdmin)
	return nil
}
