package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type StoryRepository interface {
	Create(story *model.Story) error
	FindAllByUserID(userID uint, page, limit int, status, search string) ([]model.Story, int64, error)
	FindByID(id uint, userID uint) (*model.Story, error)
	FindWithDetails(id uint, userID uint) (*model.Story, error)
	Update(story *model.Story) error
	Delete(id uint, userID uint) error
}

type storyRepository struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) StoryRepository {
	return &storyRepository{db: db}
}

func (r *storyRepository) Create(story *model.Story) error {
	return r.db.Create(story).Error
}

func (r *storyRepository) FindAllByUserID(userID uint, page, limit int, status, search string) ([]model.Story, int64, error) {
	var stories []model.Story
	var total int64

	query := r.db.Model(&model.Story{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR genre ILIKE ? OR synopsis ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&stories).Error
	if err != nil {
		return nil, 0, err
	}

	return stories, total, nil
}

func (r *storyRepository) FindByID(id uint, userID uint) (*model.Story, error) {
	var story model.Story
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&story).Error
	if err != nil {
		return nil, err
	}
	return &story, nil
}

func (r *storyRepository) FindWithDetails(id uint, userID uint) (*model.Story, error) {
	var story model.Story
	err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Categories.Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("system_category_items.order_index ASC")
		}).
		Preload("Locations").
		Preload("Factions").
		Preload("Characters").
		Preload("Plots", func(db *gorm.DB) *gorm.DB {
			return db.Order("plots.order_index ASC")
		}).
		Preload("Items").
		First(&story).Error

	if err != nil {
		return nil, err
	}
	return &story, nil
}

func (r *storyRepository) Update(story *model.Story) error {
	return r.db.Save(story).Error
}

func (r *storyRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Story{}).Error
}
