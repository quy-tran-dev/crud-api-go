package world

import (
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

type CharacterService interface {
	CreateCharacter(storyID, userID uint, req model.CreateCharacterRequest) (*CharacterCreateResult, error)
	GetCharacters(storyID, userID uint, factionID *uint, role, level string) ([]model.Character, error)
	GetCharacterByID(id, storyID, userID uint) (*model.Character, error)
	UpdateCharacter(id, storyID, userID uint, req model.UpdateCharacterRequest) (*model.Character, error)
	DeleteCharacter(id, storyID, userID uint) error
}

type characterService struct {
	characterRepo repository.CharacterRepository
	categoryRepo  repository.CategoryRepository
	storyRepo     repository.StoryRepository
}

func NewCharacterService(
	characterRepo repository.CharacterRepository,
	categoryRepo repository.CategoryRepository,
	storyRepo repository.StoryRepository,
) CharacterService {
	return &characterService{
		characterRepo: characterRepo,
		categoryRepo:  categoryRepo,
		storyRepo:     storyRepo,
	}
}

func (s *characterService) CreateCharacter(storyID, userID uint, req model.CreateCharacterRequest) (*CharacterCreateResult, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	character := &model.Character{
		StoryID:       storyID,
		Name:          strings.TrimSpace(req.Name),
		Alias:         strings.TrimSpace(req.Alias),
		Role:          strings.TrimSpace(req.Role),
		Gender:        strings.TrimSpace(req.Gender),
		Age:           strings.TrimSpace(req.Age),
		FactionID:     req.FactionID,
		FactionRank:   strings.TrimSpace(req.FactionRank),
		BirthplaceID:  req.BirthplaceID,
		CurrentLevel:  strings.TrimSpace(req.CurrentLevel),
		Abilities:     req.Abilities,
		OriginStory:   req.OriginStory,
		CharacterArc:  req.CharacterArc,
		Relationships: req.Relationships,
		Inventory:     req.Inventory,
		CategoryAttrs: req.CategoryAttrs,
		IsDraft:       true,
	}

	if err := s.characterRepo.Create(character); err != nil {
		return nil, err
	}

	// 🔍 Cơ chế RECHECK / AUTO-DISCOVERY:
	// Tự động quét xem các giá trị tác giả vừa gõ có từ nào mới chưa có trong Danh mục không
	var detectedItems []model.DetectedNewItem

	// 1. Quét Cảnh giới (CurrentLevel)
	if character.CurrentLevel != "" {
		if detected, _ := s.checkNewItem(storyID, model.CategoryCodePowerLevel, character.CurrentLevel); detected != nil {
			detectedItems = append(detectedItems, *detected)
		}
	}

	// 2. Quét các trường động trong CategoryAttrs
	for _, attr := range character.CategoryAttrs {
		if attr.Value != "" && attr.CategoryCode != "" {
			if detected, _ := s.checkNewItem(storyID, attr.CategoryCode, attr.Value); detected != nil {
				detectedItems = append(detectedItems, *detected)
			}
		}
	}

	return &CharacterCreateResult{
		Character:        character,
		DetectedNewItems: detectedItems,
	}, nil
}

func (s *characterService) checkNewItem(storyID uint, categoryCode, value string) (*model.DetectedNewItem, error) {
	val := strings.TrimSpace(value)
	category, err := s.categoryRepo.FindCategoryByCode(storyID, categoryCode)
	if err != nil {
		return &model.DetectedNewItem{
			CategoryCode: categoryCode,
			NewValue:     val,
			Message:      "Phát hiện giá trị '" + val + "' thuộc danh mục '" + categoryCode + "' chưa được khởi tạo.",
		}, nil
	}

	existingItem, err := s.categoryRepo.FindCategoryItemByName(category.ID, val)
	if err != nil || existingItem == nil {
		return &model.DetectedNewItem{
			CategoryCode: categoryCode,
			NewValue:     val,
			Message:      "Phát hiện mục mới '" + val + "' chưa có trong từ điển '" + category.Name + "'. Bạn có muốn tự động thêm vào Từ điển không?",
		}, nil
	}
	return nil, nil
}

func (s *characterService) GetCharacters(storyID, userID uint, factionID *uint, role, level string) ([]model.Character, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}
	return s.characterRepo.FindAllByStoryID(storyID, factionID, role, level)
}

func (s *characterService) GetCharacterByID(id, storyID, userID uint) (*model.Character, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	character, err := s.characterRepo.FindByID(id)
	if err != nil || character.StoryID != storyID {
		return nil, ErrCharacterNotFound
	}
	return character, nil
}

func (s *characterService) UpdateCharacter(id, storyID, userID uint, req model.UpdateCharacterRequest) (*model.Character, error) {
	character, err := s.GetCharacterByID(id, storyID, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		character.Name = strings.TrimSpace(*req.Name)
	}
	if req.Alias != nil {
		character.Alias = strings.TrimSpace(*req.Alias)
	}
	if req.Role != nil {
		character.Role = strings.TrimSpace(*req.Role)
	}
	if req.Gender != nil {
		character.Gender = strings.TrimSpace(*req.Gender)
	}
	if req.Age != nil {
		character.Age = strings.TrimSpace(*req.Age)
	}
	if req.FactionID != nil {
		character.FactionID = req.FactionID
	}
	if req.FactionRank != nil {
		character.FactionRank = strings.TrimSpace(*req.FactionRank)
	}
	if req.BirthplaceID != nil {
		character.BirthplaceID = req.BirthplaceID
	}
	if req.CurrentLevel != nil {
		character.CurrentLevel = strings.TrimSpace(*req.CurrentLevel)
	}
	if req.Abilities != nil {
		character.Abilities = *req.Abilities
	}
	if req.OriginStory != nil {
		character.OriginStory = *req.OriginStory
	}
	if req.CharacterArc != nil {
		character.CharacterArc = *req.CharacterArc
	}
	if req.Relationships != nil {
		character.Relationships = req.Relationships
	}
	if req.Inventory != nil {
		character.Inventory = req.Inventory
	}
	if req.CategoryAttrs != nil {
		character.CategoryAttrs = req.CategoryAttrs
	}
	if req.IsDraft != nil {
		character.IsDraft = *req.IsDraft
	}

	if err := s.characterRepo.Update(character); err != nil {
		return nil, err
	}
	return character, nil
}

func (s *characterService) DeleteCharacter(id, storyID, userID uint) error {
	if _, err := s.GetCharacterByID(id, storyID, userID); err != nil {
		return err
	}
	return s.characterRepo.Delete(id)
}
