package model

import "time"

// Item là Vật phẩm / Trang bị / Đạo cụ
type Item struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	StoryID          uint      `gorm:"not null;index" json:"story_id"`
	Name             string    `gorm:"size:255;not null" json:"name"`
	Type             string    `gorm:"size:100" json:"type"`  // "Vũ khí", "Pháp bảo", "Đan dược", "Manh mối"...
	Grade            string    `gorm:"size:100" json:"grade"` // Phẩm giai / Phân cấp: "Địa giai thượng phẩm", "Nhất phẩm", "Legendary"...
	Function         string    `gorm:"type:text" json:"function"`
	UsageCondition   string    `gorm:"type:text" json:"usage_condition"` // Cách dùng, điều kiện, tác dụng phụ
	OriginLocationID *uint     `gorm:"index" json:"origin_location_id"`
	OwnerCharacterID *uint     `gorm:"index" json:"owner_character_id"` // Nhân vật sở hữu
	OwnerFactionID   *uint     `gorm:"index" json:"owner_faction_id"`   // Hoặc trấn phái chi bảo của Gia tộc
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// DTOs
type CreateItemRequest struct {
	Name             string `json:"name" binding:"required,min=1,max=255"`
	Type             string `json:"type"`
	Grade            string `json:"grade"`
	Function         string `json:"function"`
	UsageCondition   string `json:"usage_condition"`
	OriginLocationID *uint  `json:"origin_location_id"`
	OwnerCharacterID *uint  `json:"owner_character_id"`
	OwnerFactionID   *uint  `json:"owner_faction_id"`
}

type UpdateItemRequest struct {
	Name             *string `json:"name" binding:"omitempty,min=1,max=255"`
	Type             *string `json:"type"`
	Grade            *string `json:"grade"`
	Function         *string `json:"function"`
	UsageCondition   *string `json:"usage_condition"`
	OriginLocationID *uint   `json:"origin_location_id"`
	OwnerCharacterID *uint   `json:"owner_character_id"`
	OwnerFactionID   *uint   `json:"owner_faction_id"`
}
