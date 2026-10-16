package model

import "time"

const (
	PlotTypeMain = "main_plot"
	PlotTypeSub  = "sub_plot"
	PlotTypeSide = "side_story"
)

// Plot là Tình tiết / Sự kiện trong kịch bản
type Plot struct {
	ID                        uint      `gorm:"primaryKey" json:"id"`
	StoryID                   uint      `gorm:"not null;index" json:"story_id"`
	Title                     string    `gorm:"size:255;not null" json:"title"`
	Type                      string    `gorm:"size:50" json:"type"` // "main_plot", "sub_plot"...
	OrderIndex                int       `gorm:"default:0" json:"order_index"` // Thứ tự mạch truyện
	LocationID                *uint     `gorm:"index" json:"location_id"`    // Diễn ra tại Map nào
	Summary                   string    `gorm:"type:text" json:"summary"`    // Diễn biến chi tiết
	InvolvedFactionIDs        []uint    `gorm:"serializer:json;type:jsonb" json:"involved_faction_ids,omitempty"`
	ParticipatingCharacterIDs []uint    `gorm:"serializer:json;type:jsonb" json:"participating_character_ids,omitempty"`
	FeaturedItemIDs           []uint    `gorm:"serializer:json;type:jsonb" json:"featured_item_ids,omitempty"`
	Consequences              string    `gorm:"type:text" json:"consequences"` // Hậu quả, biến cố sau tình tiết
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// DTOs
type CreatePlotRequest struct {
	Title                     string `json:"title" binding:"required,min=2,max=255"`
	Type                      string `json:"type"`
	OrderIndex                *int   `json:"order_index"`
	LocationID                *uint  `json:"location_id"`
	Summary                   string `json:"summary"`
	InvolvedFactionIDs        []uint `json:"involved_faction_ids"`
	ParticipatingCharacterIDs []uint `json:"participating_character_ids"`
	FeaturedItemIDs           []uint `json:"featured_item_ids"`
	Consequences              string `json:"consequences"`
}

type UpdatePlotRequest struct {
	Title                     *string `json:"title" binding:"omitempty,min=2,max=255"`
	Type                      *string `json:"type"`
	OrderIndex                *int    `json:"order_index"`
	LocationID                *uint   `json:"location_id"`
	Summary                   *string `json:"summary"`
	InvolvedFactionIDs        []uint  `json:"involved_faction_ids"`
	ParticipatingCharacterIDs []uint  `json:"participating_character_ids"`
	FeaturedItemIDs           []uint  `json:"featured_item_ids"`
	Consequences              *string `json:"consequences"`
}
