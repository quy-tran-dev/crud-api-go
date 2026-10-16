package world

import (
	"errors"
	"net/http"
	"strconv"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type StoryHandler struct {
	storyService world.StoryService
}

func NewStoryHandler(storyService world.StoryService) *StoryHandler {
	return &StoryHandler{storyService: storyService}
}

// CreateStory: POST /api/v1/stories
func (h *StoryHandler) CreateStory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req model.CreateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu tạo tác phẩm không hợp lệ", err)
		return
	}

	story, err := h.storyService.CreateStory(userID, req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo tác phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Tạo tác phẩm thành công (Đã tự động khởi tạo Từ điển mẫu)", story)
}

// GetStories: GET /api/v1/stories?page=1&limit=10&status=draft&search=keyword
func (h *StoryHandler) GetStories(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	search := c.Query("search")

	stories, total, err := h.storyService.GetStories(userID, page, limit, status, search)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách tác phẩm", err)
		return
	}

	handler.ResponsePaginated(c, http.StatusOK, "Lấy danh sách tác phẩm thành công", stories, page, limit, total)
}

// GetStoryDetail: GET /api/v1/stories/:id
func (h *StoryHandler) GetStoryDetail(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	story, err := h.storyService.GetStoryDetail(storyID, userID)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy chi tiết tác phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy chi tiết tác phẩm thành công", story)
}

// UpdateStory: PUT /api/v1/stories/:id
func (h *StoryHandler) UpdateStory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.UpdateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật không hợp lệ", err)
		return
	}

	story, err := h.storyService.UpdateStory(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi cập nhật tác phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Cập nhật tác phẩm thành công", story)
}

// DeleteStory: DELETE /api/v1/stories/:id
func (h *StoryHandler) DeleteStory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	if err := h.storyService.DeleteStory(storyID, userID); err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi xóa tác phẩm", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Xóa tác phẩm thành công", nil)
}

// PublishStory: POST /api/v1/stories/:id/publish
func (h *StoryHandler) PublishStory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	result, err := h.storyService.PublishStory(storyID, userID)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi xác thực tác phẩm", err)
		return
	}

	if !result.Success {
		// Có lỗi ràng buộc toàn vẹn
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": result.Message,
			"issues":  result.Issues,
		})
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, result.Message, result.Story)
}

// Helper hàm trích xuất UserID từ Gin Context (do AuthMiddleware set)
func getUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("userID")
	if !exists {
		handler.ResponseError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", nil)
		return 0, false
	}
	return val.(uint), true
}

func parseID(param string) (uint, error) {
	id, err := strconv.ParseUint(param, 10, 32)
	return uint(id), err
}
