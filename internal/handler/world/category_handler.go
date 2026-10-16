package world

import (
	"errors"
	"net/http"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService world.CategoryService
}

func NewCategoryHandler(categoryService world.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory: POST /api/v1/stories/:id/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu tạo danh mục không hợp lệ", err)
		return
	}

	cat, err := h.categoryService.CreateCategory(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo danh mục", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Tạo danh mục thành công", cat)
}

// GetCategories: GET /api/v1/stories/:id/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	cats, err := h.categoryService.GetCategories(storyID, userID)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách danh mục", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh mục thành công", cats)
}

// CreateCategoryItem: POST /api/v1/stories/:id/categories/:category_id/items
func (h *CategoryHandler) CreateCategoryItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	categoryID, err := parseID(c.Param("category_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID danh mục không hợp lệ", err)
		return
	}

	var req model.CreateCategoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu mục từ điển không hợp lệ", err)
		return
	}

	item, err := h.categoryService.CreateCategoryItem(storyID, userID, categoryID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể thêm mục vào từ điển", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Thêm mục vào từ điển thành công", item)
}

type RecheckRequest struct {
	CategoryCode string `json:"category_code" binding:"required"`
	Value        string `json:"value" binding:"required"`
}

// RecheckValue: POST /api/v1/stories/:id/recheck
func (h *CategoryHandler) RecheckValue(c *gin.Context) {
	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req RecheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu kiểm tra không hợp lệ", err)
		return
	}

	detected, err := h.categoryService.RecheckValue(storyID, req.CategoryCode, req.Value)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi kiểm tra từ điển", err)
		return
	}

	if detected == nil {
		handler.ResponseSuccess(c, http.StatusOK, "Giá trị đã tồn tại trong Từ điển Tác phẩm", gin.H{
			"is_new": false,
		})
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Phát hiện mục mới chưa có trong Từ điển", gin.H{
		"is_new":   true,
		"detected": detected,
	})
}
