package world

import (
	"errors"
	"net/http"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type CharacterHandler struct {
	characterService world.CharacterService
}

func NewCharacterHandler(characterService world.CharacterService) *CharacterHandler {
	return &CharacterHandler{characterService: characterService}
}

// CreateCharacter: POST /api/v1/stories/:id/characters
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu nhân vật không hợp lệ", err)
		return
	}

	result, err := h.characterService.CreateCharacter(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo nhân vật", err)
		return
	}

	// Trả về nhân vật kèm các cảnh báo phát hiện từ mới nếu có
	msg := "Tạo nhân vật thành công"
	if len(result.DetectedNewItems) > 0 {
		msg = "Tạo nhân vật thành công (Phát hiện mục mới chưa có trong Từ điển)"
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":            true,
		"message":            msg,
		"data":               result.Character,
		"detected_new_items": result.DetectedNewItems,
	})
}

// GetCharacters: GET /api/v1/stories/:id/characters?faction_id=1&role=protagonist&level=Kim Đan
func (h *CharacterHandler) GetCharacters(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var factionID *uint
	if fStr := c.Query("faction_id"); fStr != "" {
		if fID, err := parseID(fStr); err == nil {
			factionID = &fID
		}
	}

	role := c.Query("role")
	level := c.Query("level")

	characters, err := h.characterService.GetCharacters(storyID, userID, factionID, role, level)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách nhân vật", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh sách nhân vật thành công", characters)
}

// GetCharacterByID: GET /api/v1/stories/:id/characters/:character_id
func (h *CharacterHandler) GetCharacterByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	characterID, err := parseID(c.Param("character_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID nhân vật không hợp lệ", err)
		return
	}

	character, err := h.characterService.GetCharacterByID(characterID, storyID, userID)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Không tìm thấy nhân vật", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy thông tin nhân vật thành công", character)
}

// UpdateCharacter: PUT /api/v1/stories/:id/characters/:character_id
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	characterID, err := parseID(c.Param("character_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID nhân vật không hợp lệ", err)
		return
	}

	var req model.UpdateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật không hợp lệ", err)
		return
	}

	character, err := h.characterService.UpdateCharacter(characterID, storyID, userID, req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi cập nhật nhân vật", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Cập nhật nhân vật thành công", character)
}

// DeleteCharacter: DELETE /api/v1/stories/:id/characters/:character_id
func (h *CharacterHandler) DeleteCharacter(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	characterID, err := parseID(c.Param("character_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID nhân vật không hợp lệ", err)
		return
	}

	if err := h.characterService.DeleteCharacter(characterID, storyID, userID); err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi xóa nhân vật", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Xóa nhân vật thành công", nil)
}
