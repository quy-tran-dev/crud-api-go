package world

import (
	"errors"
	"net/http"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type FactionHandler struct {
	factionService world.FactionService
}

func NewFactionHandler(factionService world.FactionService) *FactionHandler {
	return &FactionHandler{factionService: factionService}
}

// CreateFaction: POST /api/v1/stories/:id/factions
func (h *FactionHandler) CreateFaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.CreateFactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu thế lực không hợp lệ", err)
		return
	}

	faction, err := h.factionService.CreateFaction(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo thế lực", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Tạo thế lực thành công", faction)
}

// GetFactions: GET /api/v1/stories/:id/factions?type=sect
func (h *FactionHandler) GetFactions(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	factionType := c.Query("type")
	factions, err := h.factionService.GetFactions(storyID, userID, factionType)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách thế lực", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh sách thế lực thành công", factions)
}

// GetFactionByID: GET /api/v1/stories/:id/factions/:faction_id
func (h *FactionHandler) GetFactionByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	factionID, err := parseID(c.Param("faction_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID thế lực không hợp lệ", err)
		return
	}

	faction, err := h.factionService.GetFactionByID(factionID, storyID, userID)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Không tìm thấy thế lực", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy thông tin thế lực thành công", faction)
}

// UpdateFaction: PUT /api/v1/stories/:id/factions/:faction_id
func (h *FactionHandler) UpdateFaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	factionID, err := parseID(c.Param("faction_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID thế lực không hợp lệ", err)
		return
	}

	var req model.UpdateFactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật thế lực không hợp lệ", err)
		return
	}

	faction, err := h.factionService.UpdateFaction(factionID, storyID, userID, req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi cập nhật thế lực", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Cập nhật thế lực thành công", faction)
}

// DeleteFaction: DELETE /api/v1/stories/:id/factions/:faction_id
func (h *FactionHandler) DeleteFaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	factionID, err := parseID(c.Param("faction_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID thế lực không hợp lệ", err)
		return
	}

	if err := h.factionService.DeleteFaction(factionID, storyID, userID); err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi xóa thế lực", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Xóa thế lực thành công", nil)
}
