package world

import (
	"errors"
	"net/http"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	itemService world.ItemService
}

func NewItemHandler(itemService world.ItemService) *ItemHandler {
	return &ItemHandler{itemService: itemService}
}

// CreateItem: POST /api/v1/stories/:id/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu vật phẩm không hợp lệ", err)
		return
	}

	item, err := h.itemService.CreateItem(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo vật phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Tạo vật phẩm thành công", item)
}

// GetItems: GET /api/v1/stories/:id/items?type=weapon&grade=Địa Giai&owner_id=1
func (h *ItemHandler) GetItems(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	itemType := c.Query("type")
	grade := c.Query("grade")

	var ownerID *uint
	if oStr := c.Query("owner_id"); oStr != "" {
		if oID, err := parseID(oStr); err == nil {
			ownerID = &oID
		}
	}

	items, err := h.itemService.GetItems(storyID, userID, itemType, grade, ownerID)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách vật phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh sách vật phẩm thành công", items)
}

// GetItemByID: GET /api/v1/stories/:id/items/:item_id
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	itemID, err := parseID(c.Param("item_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID vật phẩm không hợp lệ", err)
		return
	}

	item, err := h.itemService.GetItemByID(itemID, storyID, userID)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Không tìm thấy vật phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy thông tin vật phẩm thành công", item)
}

// UpdateItem: PUT /api/v1/stories/:id/items/:item_id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	itemID, err := parseID(c.Param("item_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID vật phẩm không hợp lệ", err)
		return
	}

	var req model.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật vật phẩm không hợp lệ", err)
		return
	}

	item, err := h.itemService.UpdateItem(itemID, storyID, userID, req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi cập nhật vật phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Cập nhật vật phẩm thành công", item)
}

// DeleteItem: DELETE /api/v1/stories/:id/items/:item_id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	itemID, err := parseID(c.Param("item_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID vật phẩm không hợp lệ", err)
		return
	}

	if err := h.itemService.DeleteItem(itemID, storyID, userID); err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi xóa vật phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Xóa vật phẩm thành công", nil)
}
