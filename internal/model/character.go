package model

import "time"

// RelationshipItem biểu diễn quan hệ giữa 2 nhân vật
type RelationshipItem struct {
	TargetCharacterID   uint   `json:"target_character_id"`
	TargetCharacterName string `json:"target_character_name,omitempty"`
	RelationType        string `json:"relation_type"` // "Kẻ thù", "Huynh đệ", "Sư đồ", "Người yêu"...
	Description         string `json:"description,omitempty"`
}

// InventoryRecord biểu diễn vật phẩm nhân vật đang nắm giữ
type InventoryRecord struct {
	ItemID           uint   `json:"item_id"`
	ItemName         string `json:"item_name,omitempty"`
	AcquiredAtPlotID *uint  `json:"acquired_at_plot_id,omitempty"`
	Note             string `json:"note,omitempty"` // "Nhặt được ở bí cảnh", "Gia truyền"
}

// Character là Nhân vật
type Character struct {
	ID            uint                     `gorm:"primaryKey" json:"id"`
	StoryID       uint                     `gorm:"not null;index" json:"story_id"`
	Name          string                   `gorm:"size:255;not null" json:"name"`
	Alias         string                   `gorm:"size:255" json:"alias"` // Biệt danh, đạo hiệu
	Role          string                   `gorm:"size:100" json:"role"`  // "Nhân vật chính", "Phản diện", "Phụ"...
	Gender        string                   `gorm:"size:20" json:"gender"`
	Age           string                   `gorm:"size:50" json:"age"`

	// Gia thế & Quê quán
	FactionID     *uint                    `gorm:"index" json:"faction_id"`
	FactionRank   string                   `gorm:"size:100" json:"faction_rank"` // "Thiếu gia", "Trưởng lão", "Đệ tử"...
	BirthplaceID  *uint                    `gorm:"index" json:"birthplace_id"`

	// Sức mạnh & Tuyến truyện
	CurrentLevel  string                   `gorm:"size:100" json:"current_level"` // Cảnh giới: "Kim Đan", "Nguyên Anh"...
	Abilities     string                   `gorm:"type:text" json:"abilities"`    // Kỹ năng, công pháp sở trường
	OriginStory   string                   `gorm:"type:text" json:"origin_story"` // Quá khứ, gia thế
	CharacterArc  string                   `gorm:"type:text" json:"character_arc"` // Tuyến phát triển tâm lý

	// Các trường Dynamic JSONB
	Relationships []RelationshipItem       `gorm:"serializer:json;type:jsonb" json:"relationships,omitempty"`
	Inventory     []InventoryRecord        `gorm:"serializer:json;type:jsonb" json:"inventory,omitempty"`
	CategoryAttrs []CategoryAttributeValue `gorm:"serializer:json;type:jsonb" json:"category_attrs,omitempty"` // Phục vụ Recheck
	IsDraft       bool                     `gorm:"default:true" json:"is_draft"`

	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

// DTOs
type CreateCharacterRequest struct {
	Name          string                   `json:"name" binding:"required,min=1,max=255"`
	Alias         string                   `json:"alias"`
	Role          string                   `json:"role"`
	Gender        string                   `json:"gender"`
	Age           string                   `json:"age"`
	FactionID     *uint                    `json:"faction_id"`
	FactionRank   string                   `json:"faction_rank"`
	BirthplaceID  *uint                    `json:"birthplace_id"`
	CurrentLevel  string                   `json:"current_level"`
	Abilities     string                   `json:"abilities"`
	OriginStory   string                   `json:"origin_story"`
	CharacterArc  string                   `json:"character_arc"`
	Relationships []RelationshipItem       `json:"relationships"`
	Inventory     []InventoryRecord        `json:"inventory"`
	CategoryAttrs []CategoryAttributeValue `json:"category_attrs"`
}

type UpdateCharacterRequest struct {
	Name          *string                  `json:"name" binding:"omitempty,min=1,max=255"`
	Alias         *string                  `json:"alias"`
	Role          *string                  `json:"role"`
	Gender        *string                  `json:"gender"`
	Age           *string                  `json:"age"`
	FactionID     *uint                    `json:"faction_id"`
	FactionRank   *string                  `json:"faction_rank"`
	BirthplaceID  *uint                    `json:"birthplace_id"`
	CurrentLevel  *string                  `json:"current_level"`
	Abilities     *string                  `json:"abilities"`
	OriginStory   *string                  `json:"origin_story"`
	CharacterArc  *string                  `json:"character_arc"`
	Relationships []RelationshipItem       `json:"relationships"`
	Inventory     []InventoryRecord        `json:"inventory"`
	CategoryAttrs []CategoryAttributeValue `json:"category_attrs"`
	IsDraft       *bool                    `json:"is_draft"`
}
