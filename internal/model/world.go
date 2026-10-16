package model

import "time"

// ============================================================================
// 1. CONSTANTS & SYSTEM CODES
// ============================================================================

const (
	StoryStatusDraft    = "draft"
	StoryStatusOfficial = "official"

	// Các mã danh mục hệ thống cốt lõi (Core System Categories)
	CategoryCodePowerLevel    = "power_level"    // Hệ thống cảnh giới / cấp bậc sức mạnh
	CategoryCodeFactionType   = "faction_type"   // Phân loại thế lực (Tông môn, Gia tộc, Triều đình...)
	CategoryCodeLocationType  = "location_type"  // Phân loại địa danh (Đại lục, Bí cảnh, Sơn môn...)
	CategoryCodeCharacterRole = "character_role" // Vai trò nhân vật (Chính, Phản diện, Phụ...)
	CategoryCodeItemType      = "item_type"      // Phân loại vật phẩm (Vũ khí, Pháp bảo, Đan dược...)
	CategoryCodeElement       = "element"        // Hệ thống nguyên tố / linh căn / trường phái
)

// ============================================================================
// 2. EMBEDDED STRUCTS PHỤ TRỢ CHO JSONB
// ============================================================================

// CategoryAttributeValue lưu trữ một giá trị động gắn với danh mục trong JSONB
// Phục vụ cơ chế Auto-Discovery / Recheck thông minh
type CategoryAttributeValue struct {
	CategoryCode string `json:"category_code"` // Ví dụ: "power_level", "element"
	Value        string `json:"value"`         // Ví dụ: "Ngưng Khí Kỳ", "Lôi Linh Căn"
	Note         string `json:"note,omitempty"`
}

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

// DiplomacyRecord biểu diễn quan hệ ngoại giao giữa 2 thế lực
type DiplomacyRecord struct {
	TargetFactionID   uint   `json:"target_faction_id"`
	TargetFactionName string `json:"target_faction_name,omitempty"`
	Status            string `json:"status"` // "ally" (Đồng minh), "enemy" (Thù địch), "neutral" (Trung lập)
	Reason            string `json:"reason,omitempty"`
}

// DetectedNewItem là DTO trả về khi tính năng Recheck phát hiện tác giả gõ từ mới
type DetectedNewItem struct {
	CategoryCode string `json:"category_code"`
	NewValue     string `json:"new_value"`
	Message      string `json:"message"`
}

// ============================================================================
// 3. MASTER DATA: BỘ TỪ ĐIỂN DANH MỤC COMMON GẮN VỚI STORY
// ============================================================================

// SystemCategory là danh mục lớn thuộc về một Tác phẩm
type SystemCategory struct {
	ID          uint                 `gorm:"primaryKey" json:"id"`
	StoryID     uint                 `gorm:"not null;index" json:"story_id"`
	Code        string               `gorm:"size:100;not null" json:"code"` // "power_level", "faction_type"...
	Name        string               `gorm:"size:255;not null" json:"name"` // "Hệ thống Cảnh Giới", "Phân Loại Thế Lực"
	Description string               `gorm:"type:text" json:"description"`
	Items       []SystemCategoryItem `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// SystemCategoryItem là mục chi tiết bên trong danh mục (Nguồn sáng tạo)
type SystemCategoryItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Code        string         `gorm:"size:100" json:"code"`                   // Mã định danh (tùy chọn)
	Name        string         `gorm:"size:255;not null" json:"name"`          // "Kim Đan Kỳ", "Tông Môn Chính Đạo"...
	OrderIndex  int            `gorm:"default:0" json:"order_index"`           // Cấp bậc cao thấp / Thứ tự
	Description string         `gorm:"type:text" json:"description"`
	Metadata    map[string]any `gorm:"serializer:json;type:jsonb" json:"metadata,omitempty"` // Các chỉ số bổ sung
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ============================================================================
// 4. CÁC THỰC THỂ WORLD-BUILDING CHÍNH (GORM ENTITIES)
// ============================================================================

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

// Location là Bản đồ / Địa danh (Hỗ trợ cấu trúc cây cha - con)
type Location struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	StoryID     uint       `gorm:"not null;index" json:"story_id"`
	ParentID    *uint      `gorm:"index" json:"parent_id"` // Cây phân cấp địa lý
	Name        string     `gorm:"size:255;not null" json:"name"`
	Type        string     `gorm:"size:100" json:"type"` // Lưu text tự nhiên: "Đại lục", "Sơn môn", "Bí cảnh"...
	DangerLevel string     `gorm:"size:50" json:"danger_level"`
	Description string     `gorm:"type:text" json:"description"`
	Children    []Location `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Faction là Thế lực / Gia tộc / Tông môn / Tổ chức
type Faction struct {
	ID             uint              `gorm:"primaryKey" json:"id"`
	StoryID        uint              `gorm:"not null;index" json:"story_id"`
	Name           string            `gorm:"size:255;not null" json:"name"`
	Type           string            `gorm:"size:100" json:"type"` // "Gia tộc", "Tông môn", "Triều đình"...
	HeadquarterID  *uint             `gorm:"index" json:"headquarter_id"` // Trụ sở đóng tại Map nào
	LeaderID       *uint             `gorm:"index" json:"leader_id"`      // Người đứng đầu
	PowerTier      string            `gorm:"size:100" json:"power_tier"`  // Cấp bậc thế lực
	Description    string            `gorm:"type:text" json:"description"`
	Diplomacy      []DiplomacyRecord `gorm:"serializer:json;type:jsonb" json:"diplomacy,omitempty"` // Mối quan hệ ngoại giao
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// Character là Nhân vật
type Character struct {
	ID             uint                     `gorm:"primaryKey" json:"id"`
	StoryID        uint                     `gorm:"not null;index" json:"story_id"`
	Name           string                   `gorm:"size:255;not null" json:"name"`
	Alias          string                   `gorm:"size:255" json:"alias"` // Biệt danh, đạo hiệu
	Role           string                   `gorm:"size:100" json:"role"`  // "Nhân vật chính", "Phản diện", "Phụ"...
	Gender         string                   `gorm:"size:20" json:"gender"`
	Age            string                   `gorm:"size:50" json:"age"`

	// Gia thế & Quê quán
	FactionID      *uint                    `gorm:"index" json:"faction_id"`
	FactionRank    string                   `gorm:"size:100" json:"faction_rank"` // "Thiếu gia", "Trưởng lão", "Đệ tử"...
	BirthplaceID   *uint                    `gorm:"index" json:"birthplace_id"`

	// Sức mạnh & Tuyến truyện
	CurrentLevel   string                   `gorm:"size:100" json:"current_level"` // Cảnh giới: "Kim Đan", "Nguyên Anh"...
	Abilities      string                   `gorm:"type:text" json:"abilities"`    // Kỹ năng, công pháp sở trường
	OriginStory    string                   `gorm:"type:text" json:"origin_story"` // Quá khứ, gia thế
	CharacterArc   string                   `gorm:"type:text" json:"character_arc"` // Tuyến phát triển tâm lý

	// Các trường Dynamic JSONB
	Relationships  []RelationshipItem       `gorm:"serializer:json;type:jsonb" json:"relationships,omitempty"`
	Inventory      []InventoryRecord        `gorm:"serializer:json;type:jsonb" json:"inventory,omitempty"`
	CategoryAttrs  []CategoryAttributeValue `gorm:"serializer:json;type:jsonb" json:"category_attrs,omitempty"` // Phục vụ Recheck
	IsDraft        bool                     `gorm:"default:true" json:"is_draft"`

	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

// Plot là Tình tiết / Sự kiện trong kịch bản
type Plot struct {
	ID                         uint      `gorm:"primaryKey" json:"id"`
	StoryID                    uint      `gorm:"not null;index" json:"story_id"`
	Title                      string    `gorm:"size:255;not null" json:"title"`
	Type                       string    `gorm:"size:50" json:"type"` // "Tình tiết chính", "Tình tiết phụ"...
	OrderIndex                 int       `gorm:"default:0" json:"order_index"` // Thứ tự mạch truyện
	LocationID                 *uint     `gorm:"index" json:"location_id"`    // Diễn ra tại Map nào
	Summary                    string    `gorm:"type:text" json:"summary"`    // Diễn biến chi tiết
	InvolvedFactionIDs         []uint    `gorm:"serializer:json;type:jsonb" json:"involved_faction_ids,omitempty"`
	ParticipatingCharacterIDs  []uint    `gorm:"serializer:json;type:jsonb" json:"participating_character_ids,omitempty"`
	FeaturedItemIDs            []uint    `gorm:"serializer:json;type:jsonb" json:"featured_item_ids,omitempty"`
	Consequences               string    `gorm:"type:text" json:"consequences"` // Hậu quả, biến cố sau tình tiết
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// Item là Vật phẩm / Trang bị / Đạo cụ
type Item struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	StoryID          uint      `gorm:"not null;index" json:"story_id"`
	Name             string    `gorm:"size:255;not null" json:"name"`
	Type             string    `gorm:"size:100" json:"type"` // "Vũ khí", "Pháp bảo", "Đan dược", "Manh mối"...
	Function         string    `gorm:"type:text" json:"function"`
	UsageCondition   string    `gorm:"type:text" json:"usage_condition"` // Cách dùng, điều kiện, tác dụng phụ
	OriginLocationID *uint     `gorm:"index" json:"origin_location_id"`
	OwnerCharacterID *uint     `gorm:"index" json:"owner_character_id"` // Nhân vật sở hữu
	OwnerFactionID   *uint     `gorm:"index" json:"owner_faction_id"`   // Hoặc trấn phái chi bảo của Gia tộc
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ============================================================================
// 5. REQUEST DTOs
// ============================================================================

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

type CreateCategoryRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=100"`
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description"`
}

type CreateCategoryItemRequest struct {
	Code        string         `json:"code"`
	Name        string         `json:"name" binding:"required,min=1,max=255"`
	OrderIndex  *int           `json:"order_index"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata"`
}

type CreateLocationRequest struct {
	ParentID    *uint  `json:"parent_id"`
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Type        string `json:"type"`
	DangerLevel string `json:"danger_level"`
	Description string `json:"description"`
}

type CreateFactionRequest struct {
	Name          string            `json:"name" binding:"required,min=2,max=255"`
	Type          string            `json:"type"`
	HeadquarterID *uint             `json:"headquarter_id"`
	LeaderID      *uint             `json:"leader_id"`
	PowerTier     string            `json:"power_tier"`
	Description   string            `json:"description"`
	Diplomacy     []DiplomacyRecord `json:"diplomacy"`
}

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

type CreateItemRequest struct {
	Name             string `json:"name" binding:"required,min=1,max=255"`
	Type             string `json:"type"`
	Function         string `json:"function"`
	UsageCondition   string `json:"usage_condition"`
	OriginLocationID *uint  `json:"origin_location_id"`
	OwnerCharacterID *uint  `json:"owner_character_id"`
	OwnerFactionID   *uint  `json:"owner_faction_id"`
}
