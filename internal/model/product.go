package model

import "time"

// Product đại diện cho bảng products trong PostgreSQL
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest là DTO nhận dữ liệu khi tạo mới sản phẩm
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       *int    `json:"stock" binding:"required,gte=0"`
}

// UpdateProductRequest là DTO nhận dữ liệu khi cập nhật sản phẩm
type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
}
