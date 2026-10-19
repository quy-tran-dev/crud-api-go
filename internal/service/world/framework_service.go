package world

import (
	"errors"
	"fmt"
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

var (
	ErrFrameworkNotFound = errors.New("không tìm thấy bộ khung thế giới mẫu")
	ErrUnauthorized      = errors.New("bạn không có quyền thực hiện thao tác này")
)

type FrameworkService interface {
	CreateFramework(userID uint, req model.CreateFrameworkRequest) (*model.FrameworkTemplate, error)
	SaveFrameworkFromStory(storyID, userID uint, req model.SaveFrameworkFromStoryRequest) (*model.FrameworkTemplate, error)
	GetMyFrameworks(userID uint) ([]model.FrameworkTemplate, error)
	GetPublicFrameworks(genre, search string) ([]model.FrameworkTemplate, error)
	GetFrameworkByID(id uint) (*model.FrameworkTemplate, error)
	UpdateFramework(id, userID uint, req model.UpdateFrameworkRequest) (*model.FrameworkTemplate, error)
	DeleteFramework(id, userID uint) error
	ApplyFrameworkToStory(storyID, frameworkID, userID uint) error
}

type frameworkService struct {
	frameworkRepo repository.FrameworkRepository
	storyRepo     repository.StoryRepository
	categoryRepo  repository.CategoryRepository
}

func NewFrameworkService(
	frameworkRepo repository.FrameworkRepository,
	storyRepo repository.StoryRepository,
	categoryRepo repository.CategoryRepository,
) FrameworkService {
	return &frameworkService{
		frameworkRepo: frameworkRepo,
		storyRepo:     storyRepo,
		categoryRepo:  categoryRepo,
	}
}

func (s *frameworkService) CreateFramework(userID uint, req model.CreateFrameworkRequest) (*model.FrameworkTemplate, error) {
	tmpl := &model.FrameworkTemplate{
		UserID:      userID,
		Name:        strings.TrimSpace(req.Name),
		Genre:       strings.TrimSpace(req.Genre),
		Description: req.Description,
		IsPublic:    req.IsPublic,
		UsageCount:  0,
		Data:        req.Data,
	}

	if err := s.frameworkRepo.Create(tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *frameworkService) SaveFrameworkFromStory(storyID, userID uint, req model.SaveFrameworkFromStoryRequest) (*model.FrameworkTemplate, error) {
	story, err := s.storyRepo.FindByID(storyID, userID)
	if err != nil {
		return nil, ErrStoryNotFound
	}

	categories, err := s.categoryRepo.FindCategoriesByStoryID(storyID)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf("tác phẩm này chưa có danh mục thiết lập nào để lưu thành bộ khung")
	}

	var templateData []model.TemplateCategory
	for _, cat := range categories {
		var catItems []model.TemplateCategoryItem
		for _, item := range cat.Items {
			catItems = append(catItems, model.TemplateCategoryItem{
				Name:        item.Name,
				OrderIndex:  item.OrderIndex,
				Description: item.Description,
			})
		}

		templateData = append(templateData, model.TemplateCategory{
			Code:        cat.Code,
			Name:        cat.Name,
			Description: cat.Description,
			Items:       catItems,
		})
	}

	tmpl := &model.FrameworkTemplate{
		UserID:      userID,
		Name:        strings.TrimSpace(req.Name),
		Genre:       story.Genre,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		UsageCount:  0,
		Data:        templateData,
	}

	if err := s.frameworkRepo.Create(tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *frameworkService) GetMyFrameworks(userID uint) ([]model.FrameworkTemplate, error) {
	return s.frameworkRepo.FindByUserID(userID)
}

func (s *frameworkService) GetPublicFrameworks(genre, search string) ([]model.FrameworkTemplate, error) {
	return s.frameworkRepo.FindPublic(genre, search)
}

func (s *frameworkService) GetFrameworkByID(id uint) (*model.FrameworkTemplate, error) {
	tmpl, err := s.frameworkRepo.FindByID(id)
	if err != nil {
		return nil, ErrFrameworkNotFound
	}
	return tmpl, nil
}

func (s *frameworkService) UpdateFramework(id, userID uint, req model.UpdateFrameworkRequest) (*model.FrameworkTemplate, error) {
	tmpl, err := s.frameworkRepo.FindByID(id)
	if err != nil {
		return nil, ErrFrameworkNotFound
	}

	if tmpl.UserID != userID {
		return nil, ErrUnauthorized
	}

	if req.Name != nil {
		tmpl.Name = strings.TrimSpace(*req.Name)
	}
	if req.Genre != nil {
		tmpl.Genre = strings.TrimSpace(*req.Genre)
	}
	if req.Description != nil {
		tmpl.Description = *req.Description
	}
	if req.IsPublic != nil {
		tmpl.IsPublic = *req.IsPublic
	}
	if req.Data != nil {
		tmpl.Data = *req.Data
	}

	if err := s.frameworkRepo.Update(tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *frameworkService) DeleteFramework(id, userID uint) error {
	tmpl, err := s.frameworkRepo.FindByID(id)
	if err != nil {
		return ErrFrameworkNotFound
	}

	if tmpl.UserID != userID {
		return ErrUnauthorized
	}

	return s.frameworkRepo.Delete(id)
}

func (s *frameworkService) ApplyFrameworkToStory(storyID, frameworkID, userID uint) error {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return ErrStoryNotFound
	}

	tmpl, err := s.frameworkRepo.FindByID(frameworkID)
	if err != nil {
		return ErrFrameworkNotFound
	}

	// Kiểm tra nếu template private thì chỉ chủ sở hữu mới được dùng
	if !tmpl.IsPublic && tmpl.UserID != userID {
		return ErrUnauthorized
	}

	for _, tCat := range tmpl.Data {
		cat := &model.SystemCategory{
			StoryID:     storyID,
			Code:        tCat.Code,
			Name:        tCat.Name,
			Description: tCat.Description,
		}
		if err := s.categoryRepo.CreateCategory(cat); err != nil {
			continue
		}

		for _, tItem := range tCat.Items {
			_ = s.categoryRepo.CreateCategoryItem(&model.SystemCategoryItem{
				CategoryID:  cat.ID,
				Name:        tItem.Name,
				OrderIndex:  tItem.OrderIndex,
				Description: tItem.Description,
			})
		}
	}

	_ = s.frameworkRepo.IncrementUsage(frameworkID)
	return nil
}
