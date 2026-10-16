package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type ItemRepository interface {
	Create(item *model.Item) error
	FindAllByStoryID(storyID uint, itemType, grade string, ownerID *uint) ([]model.Item, error)
	FindByID(id uint) (*model.Item, error)
	FindByIDs(storyID uint, ids []uint) ([]model.Item, error)
	Update(item *model.Item) error
	Delete(id uint) error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(item *model.Item) error {
	return r.db.Create(item).Error
}

func (r *itemRepository) FindAllByStoryID(storyID uint, itemType, grade string, ownerID *uint) ([]model.Item, error) {
	var items []model.Item
	query := r.db.Where("story_id = ?", storyID)

	if itemType != "" {
		query = query.Where("type = ?", itemType)
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if ownerID != nil {
		query = query.Where("owner_character_id = ?", *ownerID)
	}

	err := query.Order("id ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *itemRepository) FindByID(id uint) (*model.Item, error) {
	var item model.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) FindByIDs(storyID uint, ids []uint) ([]model.Item, error) {
	if len(ids) == 0 {
		return []model.Item{}, nil
	}
	var items []model.Item
	err := r.db.Where("story_id = ? AND id IN ?", storyID, ids).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *itemRepository) Update(item *model.Item) error {
	return r.db.Save(item).Error
}

func (r *itemRepository) Delete(id uint) error {
	return r.db.Delete(&model.Item{}, id).Error
}
