package world

import (
	"errors"
	"net/http"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationService world.LocationService
}

func NewLocationHandler(locationService world.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

// CreateLocation: POST /api/v1/stories/:id/locations
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu địa danh không hợp lệ", err)
		return
	}

	loc, err := h.locationService.CreateLocation(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo địa danh", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Tạo địa danh thành công", loc)
}

// GetLocations: GET /api/v1/stories/:id/locations?tree=true
func (h *LocationHandler) GetLocations(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	treeMode := c.Query("tree") == "true"
	locs, err := h.locationService.GetLocations(storyID, userID, treeMode)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách địa danh", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh sách địa danh thành công", locs)
}

// GetLocationByID: GET /api/v1/stories/:id/locations/:location_id
func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	locationID, err := parseID(c.Param("location_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID địa danh không hợp lệ", err)
		return
	}

	loc, err := h.locationService.GetLocationByID(locationID, storyID, userID)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Không tìm thấy địa danh", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy thông tin địa danh thành công", loc)
}

// UpdateLocation: PUT /api/v1/stories/:id/locations/:location_id
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	locationID, err := parseID(c.Param("location_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID địa danh không hợp lệ", err)
		return
	}

	var req model.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật địa danh không hợp lệ", err)
		return
	}

	loc, err := h.locationService.UpdateLocation(locationID, storyID, userID, req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi cập nhật địa danh", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Cập nhật địa danh thành công", loc)
}

// DeleteLocation: DELETE /api/v1/stories/:id/locations/:location_id
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	locationID, err := parseID(c.Param("location_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID địa danh không hợp lệ", err)
		return
	}

	if err := h.locationService.DeleteLocation(locationID, storyID, userID); err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi xóa địa danh", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Xóa địa danh thành công", nil)
}
