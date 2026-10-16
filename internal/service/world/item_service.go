package world

import (
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

type ItemService interface {
	CreateItem(storyID, userID uint, req model.CreateItemRequest) (*model.Item, error)
	GetItems(storyID, userID uint, itemType, grade string, ownerID *uint) ([]model.Item, error)
	GetItemByID(id, storyID, userID uint) (*model.Item, error)
	UpdateItem(id, storyID, userID uint, req model.UpdateItemRequest) (*model.Item, error)
	DeleteItem(id, storyID, userID uint) error
}

type itemService struct {
	itemRepo  repository.ItemRepository
	storyRepo repository.StoryRepository
}

func NewItemService(itemRepo repository.ItemRepository, storyRepo repository.StoryRepository) ItemService {
	return &itemService{
		itemRepo:  itemRepo,
		storyRepo: storyRepo,
	}
}

func (s *itemService) CreateItem(storyID, userID uint, req model.CreateItemRequest) (*model.Item, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	item := &model.Item{
		StoryID:          storyID,
		Name:             strings.TrimSpace(req.Name),
		Type:             strings.TrimSpace(req.Type),
		Grade:            strings.TrimSpace(req.Grade),
		Function:         req.Function,
		UsageCondition:   req.UsageCondition,
		OriginLocationID: req.OriginLocationID,
		OwnerCharacterID: req.OwnerCharacterID,
		OwnerFactionID:   req.OwnerFactionID,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *itemService) GetItems(storyID, userID uint, itemType, grade string, ownerID *uint) ([]model.Item, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}
	return s.itemRepo.FindAllByStoryID(storyID, itemType, grade, ownerID)
}

func (s *itemService) GetItemByID(id, storyID, userID uint) (*model.Item, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	item, err := s.itemRepo.FindByID(id)
	if err != nil || item.StoryID != storyID {
		return nil, ErrItemNotFound
	}
	return item, nil
}

func (s *itemService) UpdateItem(id, storyID, userID uint, req model.UpdateItemRequest) (*model.Item, error) {
	item, err := s.GetItemByID(id, storyID, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		item.Name = strings.TrimSpace(*req.Name)
	}
	if req.Type != nil {
		item.Type = strings.TrimSpace(*req.Type)
	}
	if req.Grade != nil {
		item.Grade = strings.TrimSpace(*req.Grade)
	}
	if req.Function != nil {
		item.Function = *req.Function
	}
	if req.UsageCondition != nil {
		item.UsageCondition = *req.UsageCondition
	}
	if req.OriginLocationID != nil {
		item.OriginLocationID = req.OriginLocationID
	}
	if req.OwnerCharacterID != nil {
		item.OwnerCharacterID = req.OwnerCharacterID
	}
	if req.OwnerFactionID != nil {
		item.OwnerFactionID = req.OwnerFactionID
	}

	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *itemService) DeleteItem(id, storyID, userID uint) error {
	if _, err := s.GetItemByID(id, storyID, userID); err != nil {
		return err
	}
	return s.itemRepo.Delete(id)
}
