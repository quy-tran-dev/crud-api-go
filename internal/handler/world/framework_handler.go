package world

import (
	"errors"
	"net/http"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type FrameworkHandler struct {
	frameworkService world.FrameworkService
}

func NewFrameworkHandler(frameworkService world.FrameworkService) *FrameworkHandler {
	return &FrameworkHandler{frameworkService: frameworkService}
}

// CreateFramework: POST /api/v1/frameworks
func (h *FrameworkHandler) CreateFramework(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req model.CreateFrameworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu bộ khung không hợp lệ", err)
		return
	}

	tmpl, err := h.frameworkService.CreateFramework(userID, req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo bộ khung thế giới", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Tạo bộ khung thế giới thành công", tmpl)
}

// SaveFrameworkFromStory: POST /api/v1/stories/:id/frameworks
func (h *FrameworkHandler) SaveFrameworkFromStory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.SaveFrameworkFromStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu lưu bộ khung không hợp lệ", err)
		return
	}

	tmpl, err := h.frameworkService.SaveFrameworkFromStory(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lưu bộ khung từ tác phẩm: "+err.Error(), err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Lưu bộ khung thế giới thành công! Bạn có thể tái sử dụng hoặc chia sẻ với cộng đồng.", tmpl)
}

// GetMyFrameworks: GET /api/v1/frameworks/my
func (h *FrameworkHandler) GetMyFrameworks(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	list, err := h.frameworkService.GetMyFrameworks(userID)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách bộ khung của bạn", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh sách bộ khung thành công", list)
}

// GetPublicFrameworks: GET /api/v1/frameworks/public?genre=Tu Tiên&search=phàm nhân
func (h *FrameworkHandler) GetPublicFrameworks(c *gin.Context) {
	genre := c.Query("genre")
	search := c.Query("search")

	list, err := h.frameworkService.GetPublicFrameworks(genre, search)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách thư viện cộng đồng", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh sách thư viện cộng đồng thành công", list)
}

// GetFrameworkByID: GET /api/v1/frameworks/:id
func (h *FrameworkHandler) GetFrameworkByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID bộ khung không hợp lệ", err)
		return
	}

	tmpl, err := h.frameworkService.GetFrameworkByID(id)
	if err != nil {
		if errors.Is(err, world.ErrFrameworkNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể lấy thông tin bộ khung", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy thông tin bộ khung thành công", tmpl)
}

// UpdateFramework: PUT /api/v1/frameworks/:id
func (h *FrameworkHandler) UpdateFramework(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	id, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID bộ khung không hợp lệ", err)
		return
	}

	var req model.UpdateFrameworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật bộ khung không hợp lệ", err)
		return
	}

	tmpl, err := h.frameworkService.UpdateFramework(id, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrFrameworkNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, world.ErrUnauthorized) {
			handler.ResponseError(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể cập nhật bộ khung", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Cập nhật bộ khung thành công", tmpl)
}

// DeleteFramework: DELETE /api/v1/frameworks/:id
func (h *FrameworkHandler) DeleteFramework(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	id, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID bộ khung không hợp lệ", err)
		return
	}

	if err := h.frameworkService.DeleteFramework(id, userID); err != nil {
		if errors.Is(err, world.ErrFrameworkNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, world.ErrUnauthorized) {
			handler.ResponseError(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể xóa bộ khung", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Xóa bộ khung thành công", nil)
}

// ApplyFrameworkToStory: POST /api/v1/stories/:id/apply-framework/:framework_id
func (h *FrameworkHandler) ApplyFrameworkToStory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	frameworkID, err := parseID(c.Param("framework_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID bộ khung không hợp lệ", err)
		return
	}

	if err := h.frameworkService.ApplyFrameworkToStory(storyID, frameworkID, userID); err != nil {
		if errors.Is(err, world.ErrStoryNotFound) || errors.Is(err, world.ErrFrameworkNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, world.ErrUnauthorized) {
			handler.ResponseError(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi áp dụng bộ khung vào tác phẩm: "+err.Error(), err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Đã áp dụng bộ khung vào kịch bản thành công!", nil)
}
