package world

import (
	"fmt"
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

type CategoryService interface {
	CreateCategory(storyID, userID uint, req model.CreateCategoryRequest) (*model.SystemCategory, error)
	GetCategories(storyID, userID uint) ([]model.SystemCategory, error)
	CreateCategoryItem(storyID, userID, categoryID uint, req model.CreateCategoryItemRequest) (*model.SystemCategoryItem, error)
	RecheckValue(storyID uint, categoryCode, value string) (*model.DetectedNewItem, error)
	DeleteCategory(id, storyID, userID uint) error
	DeleteCategoryItem(id, storyID, userID uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	storyRepo    repository.StoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, storyRepo repository.StoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		storyRepo:    storyRepo,
	}
}

func (s *categoryService) CreateCategory(storyID, userID uint, req model.CreateCategoryRequest) (*model.SystemCategory, error) {
	// Kiểm tra quyền sở hữu tác phẩm
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	category := &model.SystemCategory{
		StoryID:     storyID,
		Code:        strings.ToLower(strings.TrimSpace(req.Code)),
		Name:        strings.TrimSpace(req.Name),
		Description: req.Description,
	}

	if err := s.categoryRepo.CreateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetCategories(storyID, userID uint) ([]model.SystemCategory, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}
	return s.categoryRepo.FindCategoriesByStoryID(storyID)
}

func (s *categoryService) CreateCategoryItem(storyID, userID, categoryID uint, req model.CreateCategoryItemRequest) (*model.SystemCategoryItem, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	order := 0
	if req.OrderIndex != nil {
		order = *req.OrderIndex
	}

	code := req.Code
	if code == "" {
		code = strings.ToLower(strings.ReplaceAll(req.Name, " ", "_"))
	}

	item := &model.SystemCategoryItem{
		CategoryID:  categoryID,
		Code:        code,
		Name:        strings.TrimSpace(req.Name),
		OrderIndex:  order,
		Description: req.Description,
		Metadata:    req.Metadata,
	}

	if err := s.categoryRepo.CreateCategoryItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

// RecheckValue quét xem giá trị tác giả vừa gõ đã tồn tại trong Danh mục chưa
func (s *categoryService) RecheckValue(storyID uint, categoryCode, value string) (*model.DetectedNewItem, error) {
	val := strings.TrimSpace(value)
	if val == "" {
		return nil, nil
	}

	category, err := s.categoryRepo.FindCategoryByCode(storyID, categoryCode)
	if err != nil {
		// Danh mục chưa được khởi tạo trong truyện này
		return &model.DetectedNewItem{
			CategoryCode: categoryCode,
			NewValue:     val,
			Message:      fmt.Sprintf("Phát hiện giá trị '%s' thuộc danh mục chưa khai báo '%s'. Bạn có muốn tạo danh mục và mục này không?", val, categoryCode),
		}, nil
	}

	existingItem, err := s.categoryRepo.FindCategoryItemByName(category.ID, val)
	if err != nil || existingItem == nil {
		// Chưa có trong danh mục ➔ Báo phát hiện mới
		return &model.DetectedNewItem{
			CategoryCode: categoryCode,
			NewValue:     val,
			Message:      fmt.Sprintf("Phát hiện mục mới '%s' chưa có trong '%s'. Bạn có muốn tự động thêm vào Từ điển Tác phẩm không?", val, category.Name),
		}, nil
	}

	// Đã tồn tại trong từ điển
	return nil, nil
}

func (s *categoryService) DeleteCategory(id, storyID, userID uint) error {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil
	}
	return s.categoryRepo.DeleteCategory(id)
}

func (s *categoryService) DeleteCategoryItem(id, storyID, userID uint) error {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil
	}
	return s.categoryRepo.DeleteCategoryItem(id)
}
