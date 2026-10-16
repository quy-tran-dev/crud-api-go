package model

import "time"

// DiplomacyRecord biểu diễn quan hệ ngoại giao giữa 2 thế lực
type DiplomacyRecord struct {
	TargetFactionID   uint   `json:"target_faction_id"`
	TargetFactionName string `json:"target_faction_name,omitempty"`
	Status            string `json:"status"` // "ally" (Đồng minh), "enemy" (Thù địch), "neutral" (Trung lập)
	Reason            string `json:"reason,omitempty"`
}

// Faction là Thế lực / Gia tộc / Tông môn / Tổ chức
type Faction struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	StoryID       uint              `gorm:"not null;index" json:"story_id"`
	Name          string            `gorm:"size:255;not null" json:"name"`
	Type          string            `gorm:"size:100" json:"type"` // "Gia tộc", "Tông môn", "Triều đình"...
	HeadquarterID *uint             `gorm:"index" json:"headquarter_id"` // Trụ sở đóng tại Map nào
	LeaderID      *uint             `gorm:"index" json:"leader_id"`      // Người đứng đầu
	PowerTier     string            `gorm:"size:100" json:"power_tier"`  // Cấp bậc thế lực
	Description   string            `gorm:"type:text" json:"description"`
	Diplomacy     []DiplomacyRecord `gorm:"serializer:json;type:jsonb" json:"diplomacy,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// DTOs
type CreateFactionRequest struct {
	Name          string            `json:"name" binding:"required,min=2,max=255"`
	Type          string            `json:"type"`
	HeadquarterID *uint             `json:"headquarter_id"`
	LeaderID      *uint             `json:"leader_id"`
	PowerTier     string            `json:"power_tier"`
	Description   string            `json:"description"`
	Diplomacy     []DiplomacyRecord `json:"diplomacy"`
}

type UpdateFactionRequest struct {
	Name          *string           `json:"name" binding:"omitempty,min=2,max=255"`
	Type          *string           `json:"type"`
	HeadquarterID *uint             `json:"headquarter_id"`
	LeaderID      *uint             `json:"leader_id"`
	PowerTier     *string           `json:"power_tier"`
	Description   *string           `json:"description"`
	Diplomacy     []DiplomacyRecord `json:"diplomacy"`
}
