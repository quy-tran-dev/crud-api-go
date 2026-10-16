package world

import (
	"fmt"
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

type StoryService interface {
	CreateStory(userID uint, req model.CreateStoryRequest) (*model.Story, error)
	GetStories(userID uint, page, limit int, status, search string) ([]model.Story, int64, error)
	GetStoryDetail(id, userID uint) (*model.Story, error)
	UpdateStory(id, userID uint, req model.UpdateStoryRequest) (*model.Story, error)
	DeleteStory(id, userID uint) error
	PublishStory(id, userID uint) (*PublishResult, error)
}

type storyService struct {
	storyRepo     repository.StoryRepository
	categoryRepo  repository.CategoryRepository
	locationRepo  repository.LocationRepository
	factionRepo   repository.FactionRepository
	characterRepo repository.CharacterRepository
	plotRepo      repository.PlotRepository
	itemRepo      repository.ItemRepository
}

func NewStoryService(
	storyRepo repository.StoryRepository,
	categoryRepo repository.CategoryRepository,
	locationRepo repository.LocationRepository,
	factionRepo repository.FactionRepository,
	characterRepo repository.CharacterRepository,
	plotRepo repository.PlotRepository,
	itemRepo repository.ItemRepository,
) StoryService {
	return &storyService{
		storyRepo:     storyRepo,
		categoryRepo:  categoryRepo,
		locationRepo:  locationRepo,
		factionRepo:   factionRepo,
		characterRepo: characterRepo,
		plotRepo:      plotRepo,
		itemRepo:      itemRepo,
	}
}

func (s *storyService) CreateStory(userID uint, req model.CreateStoryRequest) (*model.Story, error) {
	story := &model.Story{
		UserID:       userID,
		Title:        strings.TrimSpace(req.Title),
		Genre:        strings.TrimSpace(req.Genre),
		Synopsis:     req.Synopsis,
		WorldSetting: req.WorldSetting,
		Status:       model.StoryStatusDraft,
		IsDraft:      true,
	}

	if err := s.storyRepo.Create(story); err != nil {
		return nil, err
	}

	// Tự động Seed (khởi tạo mẫu) bộ Danh mục Từ điển theo thể loại
	s.seedDefaultCategories(story.ID, story.Genre)

	return story, nil
}

func (s *storyService) GetStories(userID uint, page, limit int, status, search string) ([]model.Story, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.storyRepo.FindAllByUserID(userID, page, limit, status, search)
}

func (s *storyService) GetStoryDetail(id, userID uint) (*model.Story, error) {
	story, err := s.storyRepo.FindWithDetails(id, userID)
	if err != nil {
		return nil, ErrStoryNotFound
	}
	return story, nil
}

func (s *storyService) UpdateStory(id, userID uint, req model.UpdateStoryRequest) (*model.Story, error) {
	story, err := s.storyRepo.FindByID(id, userID)
	if err != nil {
		return nil, ErrStoryNotFound
	}

	if req.Title != nil {
		story.Title = strings.TrimSpace(*req.Title)
	}
	if req.Genre != nil {
		story.Genre = strings.TrimSpace(*req.Genre)
	}
	if req.Synopsis != nil {
		story.Synopsis = *req.Synopsis
	}
	if req.WorldSetting != nil {
		story.WorldSetting = *req.WorldSetting
	}
	if req.Status != nil {
		story.Status = *req.Status
		story.IsDraft = (*req.Status == model.StoryStatusDraft)
	}

	if err := s.storyRepo.Update(story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *storyService) DeleteStory(id, userID uint) error {
	if _, err := s.storyRepo.FindByID(id, userID); err != nil {
		return ErrStoryNotFound
	}
	return s.storyRepo.Delete(id, userID)
}

// PublishStory: Pre-publish Integrity Validator rà soát toàn bộ thực thể khi chuyển sang Official
func (s *storyService) PublishStory(id, userID uint) (*PublishResult, error) {
	story, err := s.storyRepo.FindWithDetails(id, userID)
	if err != nil {
		return nil, ErrStoryNotFound
	}

	var issues []string

	// Tạo các map để tra cứu nhanh sự tồn tại của ID
	locationMap := make(map[uint]bool)
	for _, l := range story.Locations {
		locationMap[l.ID] = true
	}

	factionMap := make(map[uint]bool)
	for _, f := range story.Factions {
		factionMap[f.ID] = true
	}

	characterMap := make(map[uint]bool)
	for _, c := range story.Characters {
		characterMap[c.ID] = true
	}

	// 1. Rà soát Nhân vật
	for _, c := range story.Characters {
		if c.FactionID != nil && !factionMap[*c.FactionID] {
			issues = append(issues, fmt.Sprintf("Nhân vật '%s' (ID: %d) đang gán Thế lực ID %d không tồn tại trong tác phẩm.", c.Name, c.ID, *c.FactionID))
		}
		if c.BirthplaceID != nil && !locationMap[*c.BirthplaceID] {
			issues = append(issues, fmt.Sprintf("Nhân vật '%s' (ID: %d) có Nơi sinh (Địa danh ID: %d) không tồn tại.", c.Name, c.ID, *c.BirthplaceID))
		}
	}

	// 2. Rà soát Tình tiết
	for _, p := range story.Plots {
		if p.LocationID != nil && !locationMap[*p.LocationID] {
			issues = append(issues, fmt.Sprintf("Tình tiết '%s' (ID: %d) diễn ra tại Địa danh ID %d không tồn tại.", p.Title, p.ID, *p.LocationID))
		}
		for _, charID := range p.ParticipatingCharacterIDs {
			if !characterMap[charID] {
				issues = append(issues, fmt.Sprintf("Tình tiết '%s' (ID: %d) có nhân vật tham gia (ID: %d) không tồn tại.", p.Title, p.ID, charID))
			}
		}
		for _, facID := range p.InvolvedFactionIDs {
			if !factionMap[facID] {
				issues = append(issues, fmt.Sprintf("Tình tiết '%s' (ID: %d) có Thế lực liên can (ID: %d) không tồn tại.", p.Title, p.ID, facID))
			}
		}
	}

	// 3. Rà soát Vật phẩm
	for _, it := range story.Items {
		if it.OwnerCharacterID != nil && !characterMap[*it.OwnerCharacterID] {
			issues = append(issues, fmt.Sprintf("Vật phẩm '%s' (ID: %d) có chủ sở hữu (Nhân vật ID: %d) không tồn tại.", it.Name, it.ID, *it.OwnerCharacterID))
		}
		if it.OwnerFactionID != nil && !factionMap[*it.OwnerFactionID] {
			issues = append(issues, fmt.Sprintf("Vật phẩm '%s' (ID: %d) có chủ quản (Thế lực ID: %d) không tồn tại.", it.Name, it.ID, *it.OwnerFactionID))
		}
		if it.OriginLocationID != nil && !locationMap[*it.OriginLocationID] {
			issues = append(issues, fmt.Sprintf("Vật phẩm '%s' (ID: %d) có Nguồn gốc (Địa danh ID: %d) không tồn tại.", it.Name, it.ID, *it.OriginLocationID))
		}
	}

	// Nếu có bất kỳ vi phạm ràng buộc nào ➔ Không cho Publish
	if len(issues) > 0 {
		return &PublishResult{
			Success: false,
			Message: fmt.Sprintf("Không thể chuyển sang Bản Chính Thức: phát hiện %d lỗi ràng buộc cần hoàn thiện.", len(issues)),
			Issues:  issues,
			Story:   story,
		}, nil
	}

	// Nếu hoàn toàn hợp lệ ➔ Kích hoạt trạng thái Official
	story.Status = model.StoryStatusOfficial
	story.IsDraft = false
	if err := s.storyRepo.Update(story); err != nil {
		return nil, err
	}

	return &PublishResult{
		Success: true,
		Message: "Chúc mừng! Tác phẩm đã được xác thực toàn vẹn và chuyển sang Bản Chính Thức thành công.",
		Story:   story,
	}, nil
}

// seedDefaultCategories tự động khởi tạo danh mục mẫu theo thể loại truyện
func (s *storyService) seedDefaultCategories(storyID uint, genre string) {
	// Danh mục Phẩm giai Vật phẩm (Chung)
	gradeCat := &model.SystemCategory{
		StoryID:     storyID,
		Code:        model.CategoryCodeItemGrade,
		Name:        "Phẩm Cấp / Đánh Giá Vật Phẩm",
		Description: "Thang đo phân cấp phẩm chất vũ khí, pháp bảo, đan dược",
	}
	_ = s.categoryRepo.CreateCategory(gradeCat)

	grades := []string{"Hoàng Giai", "Huyền Giai", "Địa Giai", "Thiên Giai", "Thần Giai"}
	for i, g := range grades {
		_ = s.categoryRepo.CreateCategoryItem(&model.SystemCategoryItem{
			CategoryID: gradeCat.ID,
			Name:       g,
			OrderIndex: i + 1,
		})
	}

	// Nếu thể loại là Tu Tiên hoặc Huyền Huyễn
	if strings.Contains(strings.ToLower(genre), "tu tiên") || strings.Contains(strings.ToLower(genre), "tiên hiệp") {
		// 1. Cảnh giới
		levelCat := &model.SystemCategory{
			StoryID:     storyID,
			Code:        model.CategoryCodePowerLevel,
			Name:        "Hệ Thống Cảnh Giới Tu Vi",
			Description: "Các cấp bậc tu luyện từ thấp đến cao",
		}
		_ = s.categoryRepo.CreateCategory(levelCat)
		levels := []string{"Luyện Khí", "Trúc Cơ", "Kim Đan", "Nguyên Anh", "Hóa Thần", "Luyện Hư"}
		for i, l := range levels {
			_ = s.categoryRepo.CreateCategoryItem(&model.SystemCategoryItem{
				CategoryID: levelCat.ID,
				Name:       l,
				OrderIndex: i + 1,
			})
		}

		// 2. Loại thế lực
		facCat := &model.SystemCategory{
			StoryID:     storyID,
			Code:        model.CategoryCodeFactionType,
			Name:        "Phân Loại Thế Lực",
			Description: "Tông phái, gia tộc, giáo phái",
		}
		_ = s.categoryRepo.CreateCategory(facCat)
		facTypes := []string{"Tông Môn Chính Đạo", "Ma Đạo Giáo Phái", "Tu Tiên Thế Gia", "Tán Tu Liên Minh"}
		for i, f := range facTypes {
			_ = s.categoryRepo.CreateCategoryItem(&model.SystemCategoryItem{
				CategoryID: facCat.ID,
				Name:       f,
				OrderIndex: i + 1,
			})
		}
	}
}
