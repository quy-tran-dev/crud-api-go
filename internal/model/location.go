package model

import "time"

// Location là Bản đồ / Địa danh (Hỗ trợ cấu trúc cây cha - con)
type Location struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	StoryID     uint       `gorm:"not null;index" json:"story_id"`
	ParentID    *uint      `gorm:"index" json:"parent_id"` // Cây phân cấp địa lý
	Name        string     `gorm:"size:255;not null" json:"name"`
	Type        string     `gorm:"size:100" json:"type"` // Text tự nhiên: "Đại lục", "Sơn môn", "Bí cảnh"...
	DangerLevel string     `gorm:"size:50" json:"danger_level"`
	Description string     `gorm:"type:text" json:"description"`
	Children    []Location `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DTOs
type CreateLocationRequest struct {
	ParentID    *uint  `json:"parent_id"`
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Type        string `json:"type"`
	DangerLevel string `json:"danger_level"`
	Description string `json:"description"`
}

type UpdateLocationRequest struct {
	ParentID    *uint   `json:"parent_id"`
	Name        *string `json:"name" binding:"omitempty,min=2,max=255"`
	Type        *string `json:"type"`
	DangerLevel *string `json:"danger_level"`
	Description *string `json:"description"`
}
