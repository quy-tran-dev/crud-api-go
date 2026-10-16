package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type LocationRepository interface {
	Create(location *model.Location) error
	FindAllByStoryID(storyID uint) ([]model.Location, error)
	FindTreeByStoryID(storyID uint) ([]model.Location, error)
	FindByID(id uint) (*model.Location, error)
	Update(location *model.Location) error
	Delete(id uint) error
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(location *model.Location) error {
	return r.db.Create(location).Error
}

func (r *locationRepository) FindAllByStoryID(storyID uint) ([]model.Location, error) {
	var locations []model.Location
	err := r.db.Where("story_id = ?", storyID).Order("id ASC").Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *locationRepository) FindTreeByStoryID(storyID uint) ([]model.Location, error) {
	var roots []model.Location
	err := r.db.Where("story_id = ? AND parent_id IS NULL", storyID).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Children")
		}).
		Find(&roots).Error
	if err != nil {
		return nil, err
	}
	return roots, nil
}

func (r *locationRepository) FindByID(id uint) (*model.Location, error) {
	var location model.Location
	err := r.db.Preload("Children").First(&location, id).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) Update(location *model.Location) error {
	return r.db.Save(location).Error
}

func (r *locationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Location{}, id).Error
}
