package model

import "time"

const (
	// Các mã danh mục hệ thống cốt lõi (Core System Categories)
	CategoryCodePowerLevel    = "power_level"    // Hệ thống cảnh giới / cấp bậc sức mạnh
	CategoryCodeFactionType   = "faction_type"   // Phân loại thế lực (Tông môn, Gia tộc, Triều đình...)
	CategoryCodeLocationType  = "location_type"  // Phân loại địa danh (Đại lục, Bí cảnh, Sơn môn...)
	CategoryCodeCharacterRole = "character_role" // Vai trò nhân vật (Chính, Phản diện, Phụ...)
	CategoryCodeItemType      = "item_type"      // Phân loại vật phẩm (Vũ khí, Pháp bảo, Đan dược...)
	CategoryCodeElement       = "element"        // Hệ thống nguyên tố / linh căn / trường phái
)

// CategoryAttributeValue lưu trữ một giá trị động gắn với danh mục trong JSONB
// Phục vụ cơ chế Auto-Discovery / Recheck thông minh
type CategoryAttributeValue struct {
	CategoryCode string `json:"category_code"` // Ví dụ: "power_level", "element"
	Value        string `json:"value"`         // Ví dụ: "Ngưng Khí Kỳ", "Lôi Linh Căn"
	Note         string `json:"note,omitempty"`
}

// DetectedNewItem là DTO trả về khi tính năng Recheck phát hiện tác giả gõ từ mới
type DetectedNewItem struct {
	CategoryCode string `json:"category_code"`
	NewValue     string `json:"new_value"`
	Message      string `json:"message"`
}

// SystemCategory là danh mục lớn thuộc về một Tác phẩm (Master Data)
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
	Code        string         `gorm:"size:100" json:"code"`
	Name        string         `gorm:"size:255;not null" json:"name"` // "Kim Đan Kỳ", "Tông Môn Chính Đạo"...
	OrderIndex  int            `gorm:"default:0" json:"order_index"`  // Thứ tự cấp bậc cao thấp
	Description string         `gorm:"type:text" json:"description"`
	Metadata    map[string]any `gorm:"serializer:json;type:jsonb" json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// DTOs
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
