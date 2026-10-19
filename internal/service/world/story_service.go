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
	gLower := strings.ToLower(strings.TrimSpace(genre))

	helperCreate := func(code, name, desc string, items []string) {
		cat := &model.SystemCategory{
			StoryID:     storyID,
			Code:        code,
			Name:        name,
			Description: desc,
		}
		if err := s.categoryRepo.CreateCategory(cat); err == nil && cat.ID > 0 {
			for i, itName := range items {
				_ = s.categoryRepo.CreateCategoryItem(&model.SystemCategoryItem{
					CategoryID:  cat.ID,
					Name:        itName,
					OrderIndex:  i + 1,
					Description: fmt.Sprintf("Mục định nghĩa chuẩn của %s", name),
				})
			}
		}
	}

	if strings.Contains(gLower, "tu tiên") || strings.Contains(gLower, "tiên hiệp") {
		// 1. Cảnh giới
		helperCreate(
			model.CategoryCodePowerLevel,
			"Hệ Thống Cảnh Giới Tu Vi",
			"Các cấp bậc tu luyện từ thấp đến cao trong thế giới Tiên Hiệp",
			[]string{"Luyện Khí Kỳ", "Trúc Cơ Kỳ", "Kim Đan Kỳ", "Nguyên Anh Kỳ", "Hóa Thần Kỳ", "Luyện Hư Kỳ", "Hợp Thể Kỳ", "Đại Thừa Kỳ", "Độ Kiếp Kỳ"},
		)
		// 2. Phân loại thế lực
		helperCreate(
			model.CategoryCodeFactionType,
			"Phân Loại Thế Lực & Môn Phái",
			"Tông môn, gia tộc, giáo phái, liên minh tu sĩ",
			[]string{"Tông Môn Chính Đạo", "Ma Đạo Giáo Phái", "Tu Tiên Thế Gia", "Tán Tu Liên Minh", "Thương Hội Tu Giới"},
		)
		// 3. Phẩm cấp vật phẩm
		helperCreate(
			model.CategoryCodeItemGrade,
			"Phẩm Cấp / Giai Vị Pháp Bảo",
			"Hệ thống phân cấp bảo vật tu tiên",
			[]string{"Phàm Giai", "Hoàng Giai", "Huyền Giai", "Địa Giai", "Thiên Giai", "Tiên Giai", "Thần Khí Thượng Cổ"},
		)
		// 4. Loại vật phẩm
		helperCreate(
			model.CategoryCodeItemType,
			"Loại Trang Bị & Vật Phẩm",
			"Phân loại các vật phẩm trong kịch bản",
			[]string{"Vũ khí", "Pháp bảo", "Đan dược", "Linh thảo", "Phù lục", "Trận bàn", "Công pháp"},
		)
	} else if strings.Contains(gLower, "cyber") || strings.Contains(gLower, "khoa học") || strings.Contains(gLower, "viễn tưởng") {
		// 1. Cyberware / Thân thể nâng cấp
		helperCreate(
			model.CategoryCodePowerLevel,
			"Cấp Bậc Nâng Cấp (Cyberware Tier)",
			"Mức độ tích hợp công nghệ và độ chịu tải thần kinh của cơ thể",
			[]string{"Tier 1 (Chợ Đen)", "Tier 2 (Dân Sự Chuẩn)", "Tier 3 (Quân Đội Tăng Cường)", "Tier 4 (Nguyên Mẫu Tập Đoàn)", "Tier 5 (Huyền Thoại Cyber)"},
		)
		// 2. Phân loại thế lực
		helperCreate(
			model.CategoryCodeFactionType,
			"Phân Loại Thế Lực & Tập Đoàn",
			"Tập đoàn khổng lồ, băng đảng đường phố, tổ chức ngầm",
			[]string{"Đại Tập Đoàn (Megacorp)", "Băng Đảng Đường Phố", "Tổ Chức Hacker Tự Do", "Lực Lượng An Ninh (Police)", "Thợ Săn Tiền Thưởng (Mercs)"},
		)
		// 3. Phẩm cấp trang bị
		helperCreate(
			model.CategoryCodeItemGrade,
			"Phẩm Cấp Công Nghệ",
			"Chất lượng và độ hiếm của vũ khí công nghệ",
			[]string{"Phổ Thông (Common)", "Cải Tiến (Custom)", "Quân Sự (Military-Grade)", "Công Nghệ Cao (High-Tech)", "Nguyên Mẫu Độc Bản (Prototype)"},
		)
		// 4. Loại trang bị
		helperCreate(
			model.CategoryCodeItemType,
			"Loại Trang Bị & Vật Phẩm",
			"Phân loại vũ khí và đồ chơi công nghệ",
			[]string{"Vũ khí thông minh", "Cấy ghép thần kinh", "Chip dữ liệu", "Phần mềm hack", "Bộ giáp tàng hình", "Thuốc kích thích"},
		)
	} else if strings.Contains(gLower, "trinh thám") || strings.Contains(gLower, "hình sự") || strings.Contains(gLower, "bắt án") {
		// 1. Cấp bậc điều tra
		helperCreate(
			model.CategoryCodePowerLevel,
			"Cấp Bậc Điều Tra & Thẩm Quyền",
			"Thẩm quyền và cấp bậc trong hệ thống cảnh sát và thám tử",
			[]string{"Thực Tập Sinh", "Điều Tra Viên", "Thanh Tra Trưởng", "Cục Trưởng Điều Tra", "Thám Tử Cố Vấn Đặc Biệt"},
		)
		// 2. Phân loại thế lực
		helperCreate(
			model.CategoryCodeFactionType,
			"Phân Loại Tổ Chức & Băng Nhóm",
			"Cơ quan hành pháp, băng nhóm tội phạm, văn phòng thám tử",
			[]string{"Sở Cảnh Sát Thành Phố", "Văn Phòng Thám Tử Tư", "Tổ Chức Tội Phạm Ngầm", "Viện Pháp Y Quốc Gia", "Viện Kiểm Sát"},
		)
		// 3. Phẩm cấp manh mối
		helperCreate(
			model.CategoryCodeItemGrade,
			"Mức Độ Quan Trọng Của Manh Mối",
			"Giá trị chứng cứ đối với vụ án",
			[]string{"Manh Mối Sơ Cấp", "Chứng Cứ Đáng Ngờ", "Manh Mối Quan Trọng", "Bằng Chứng Then Chốt", "Vật Chứng Tuyệt Mật"},
		)
		// 4. Loại vật phẩm
		helperCreate(
			model.CategoryCodeItemType,
			"Loại Trang Bị & Vật Phẩm",
			"Manh mối, thiết bị điều tra, hồ sơ vụ án",
			[]string{"Manh mối hiện trường", "Báo cáo pháp y", "Thiết bị ghi âm/nghe lén", "Vũ khí tự vệ", "Hồ sơ lưu trữ mật"},
		)
	} else if strings.Contains(gLower, "huyền huyễn") || strings.Contains(gLower, "fantasy") || strings.Contains(gLower, "ma pháp") {
		// 1. Cấp bậc ma pháp
		helperCreate(
			model.CategoryCodePowerLevel,
			"Hệ Thống Ma Pháp / Chức Nghiệp",
			"Cấp bậc sức mạnh của pháp sư và hiệp sĩ",
			[]string{"Học Đồ", "Sơ Cấp (Sắt)", "Trung Cấp (Đồng)", "Cao Cấp (Bạc)", "Đại Pháp Sư (Vàng)", "Thánh Vực", "Thần Thoại"},
		)
		// 2. Phân loại thế lực
		helperCreate(
			model.CategoryCodeFactionType,
			"Chủng Tộc & Thế Lực",
			"Vương quốc, công hội, giáo đình, gia tộc cổ",
			[]string{"Đế Quốc Nhân Tộc", "Vương Quốc Tinh Linh", "Bộ Lạc Thú Tộc", "Ma Pháp Công Hội", "Giáo Đình Thánh Quang", "Quân Đoàn Hắc Ám"},
		)
		// 3. Phẩm cấp trang bị
		helperCreate(
			model.CategoryCodeItemGrade,
			"Phẩm Cấp Bảo Vật Ma Pháp",
			"Hệ thống phân cấp bảo vật fantasy",
			[]string{"Bạch Thiết", "Thanh Đồng", "Bạch Ngân", "Hoàng Kim", "Sử Thi (Epic)", "Truyền Thuyết (Legendary)", "Thần Khí (Artifact)"},
		)
		// 4. Loại vật phẩm
		helperCreate(
			model.CategoryCodeItemType,
			"Loại Trang Bị & Vật Phẩm",
			"Phân loại trang bị fantasy",
			[]string{"Trượng ma pháp", "Vũ khí cận chiến", "Áo choàng phép", "Ma dược", "Cuộn giấy phép thuật", "Ma tinh"},
		)
	} else {
		// Thể loại Tùy chỉnh (Custom) - Cung cấp bộ khung nền tảng linh hoạt
		helperCreate(
			model.CategoryCodePowerLevel,
			"Hệ Thống Cấp Bậc / Năng Lực",
			"Thang đo cấp bậc sức mạnh tùy chỉnh",
			[]string{"Cấp F (Nhập Môn)", "Cấp E", "Cấp D", "Cấp C", "Cấp B", "Cấp A", "Cấp S (Đỉnh Cao)", "Cấp Siêu Việt"},
		)
		helperCreate(
			model.CategoryCodeFactionType,
			"Phân Loại Thế Lực & Tổ Chức",
			"Các tổ chức, liên minh trong tác phẩm",
			[]string{"Tổ Chức Chính Quyền", "Liên Minh Độc Lập", "Phe Nổi Dậy", "Tổ Chức Bí Ẩn"},
		)
		helperCreate(
			model.CategoryCodeItemGrade,
			"Phẩm Cấp Đánh Giá Trang Bị",
			"Cấp độ phẩm chất vật phẩm",
			[]string{"Cấp 1 (Cơ Bản)", "Cấp 2 (Nâng Cao)", "Cấp 3 (Hiếm)", "Cấp 4 (Cực Phẩm)", "Cấp 5 (Thần Thoại)"},
		)
		helperCreate(
			model.CategoryCodeItemType,
			"Loại Trang Bị & Vật Phẩm",
			"Phân loại các vật phẩm trong kịch bản",
			[]string{"Vũ khí", "Giáp bảo hộ", "Dụng cụ hỗ trợ", "Đạo cụ đặc biệt", "Manh mối"},
		)
	}
}
