package model

import "time"

// FrameworkTemplate là bộ khung thiết lập thế giới (cảnh giới, phẩm cấp, thế lực) có thể tái sử dụng và chia sẻ công khai
type FrameworkTemplate struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	UserID      uint               `gorm:"not null;index" json:"user_id"`
	Author      *User              `gorm:"foreignKey:UserID" json:"author,omitempty"`
	Name        string             `gorm:"size:255;not null" json:"name"`
	Genre       string             `gorm:"size:100;not null;index" json:"genre"`
	Description string             `gorm:"type:text" json:"description"`
	IsPublic    bool               `gorm:"default:false;index" json:"is_public"`
	UsageCount  int                `gorm:"default:0" json:"usage_count"`
	Data        []TemplateCategory `gorm:"serializer:json;type:jsonb" json:"data"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type TemplateCategory struct {
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Items       []TemplateCategoryItem `json:"items"`
}

type TemplateCategoryItem struct {
	Name        string `json:"name"`
	OrderIndex  int    `json:"order_index"`
	Description string `json:"description"`
}

// DTOs
type CreateFrameworkRequest struct {
	Name        string             `json:"name" binding:"required,min=2,max=255"`
	Genre       string             `json:"genre" binding:"required,min=1,max=100"`
	Description string             `json:"description"`
	IsPublic    bool               `json:"is_public"`
	Data        []TemplateCategory `json:"data" binding:"required"`
}

type UpdateFrameworkRequest struct {
	Name        *string             `json:"name" binding:"omitempty,min=2,max=255"`
	Genre       *string             `json:"genre" binding:"omitempty,min=1,max=100"`
	Description *string             `json:"description"`
	IsPublic    *bool               `json:"is_public"`
	Data        *[]TemplateCategory `json:"data"`
}

type SaveFrameworkFromStoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}
