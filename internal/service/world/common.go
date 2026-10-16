package world

import (
	"errors"

	"crud-api/internal/model"
)

var (
	ErrStoryNotFound     = errors.New("không tìm thấy tác phẩm hoặc bạn không có quyền truy cập")
	ErrCategoryNotFound  = errors.New("không tìm thấy danh mục")
	ErrLocationNotFound  = errors.New("không tìm thấy địa danh")
	ErrFactionNotFound   = errors.New("không tìm thấy thế lực")
	ErrCharacterNotFound = errors.New("không tìm thấy nhân vật")
	ErrPlotNotFound      = errors.New("không tìm thấy tình tiết")
	ErrItemNotFound      = errors.New("không tìm thấy vật phẩm")
)

// PublishResult là kết quả kiểm tra toàn vẹn khi chuyển từ Draft sang Official
type PublishResult struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Issues  []string     `json:"issues,omitempty"`
	Story   *model.Story `json:"story,omitempty"`
}

// CharacterCreateResult trả về thông tin nhân vật kèm các cảnh báo phát hiện từ mới (Recheck)
type CharacterCreateResult struct {
	Character        *model.Character        `json:"character"`
	DetectedNewItems []model.DetectedNewItem `json:"detected_new_items,omitempty"`
}
