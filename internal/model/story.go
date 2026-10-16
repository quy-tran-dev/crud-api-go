package model

import "time"

const (
	StoryStatusDraft    = "draft"
	StoryStatusOfficial = "official"
)

// Story đại diện cho Tác phẩm / Dự án kịch bản (Thư mục gốc)
type Story struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	UserID       uint             `gorm:"not null;index" json:"user_id"` // Tác giả sở hữu
	Title        string           `gorm:"size:255;not null" json:"title"`
	Genre        string           `gorm:"size:100" json:"genre"` // "Tu Tiên", "Cyberpunk", "Trinh Thám"...
	Synopsis     string           `gorm:"type:text" json:"synopsis"`
	WorldSetting string           `gorm:"type:text" json:"world_setting"` // Bối cảnh thế giới, quy tắc không gian
	Status       string           `gorm:"size:20;default:'draft'" json:"status"` // "draft" hoặc "official"
	IsDraft      bool             `gorm:"default:true" json:"is_draft"`

	// Các folder con bên trong tác phẩm
	Categories   []SystemCategory `gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE" json:"categories,omitempty"`
	Locations    []Location       `gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE" json:"locations,omitempty"`
	Factions     []Faction        `gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE" json:"factions,omitempty"`
	Characters   []Character      `gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE" json:"characters,omitempty"`
	Plots        []Plot           `gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE" json:"plots,omitempty"`
	Items        []Item           `gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE" json:"items,omitempty"`

	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// DTOs
type CreateStoryRequest struct {
	Title        string `json:"title" binding:"required,min=2,max=255"`
	Genre        string `json:"genre"`
	Synopsis     string `json:"synopsis"`
	WorldSetting string `json:"world_setting"`
}

type UpdateStoryRequest struct {
	Title        *string `json:"title" binding:"omitempty,min=2,max=255"`
	Genre        *string `json:"genre"`
	Synopsis     *string `json:"synopsis"`
	WorldSetting *string `json:"world_setting"`
	Status       *string `json:"status" binding:"omitempty,oneof=draft official"`
}
