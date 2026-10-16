package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type FactionRepository interface {
	Create(faction *model.Faction) error
	FindAllByStoryID(storyID uint, factionType string) ([]model.Faction, error)
	FindByID(id uint) (*model.Faction, error)
	Update(faction *model.Faction) error
	Delete(id uint) error
}

type factionRepository struct {
	db *gorm.DB
}

func NewFactionRepository(db *gorm.DB) FactionRepository {
	return &factionRepository{db: db}
}

func (r *factionRepository) Create(faction *model.Faction) error {
	return r.db.Create(faction).Error
}

func (r *factionRepository) FindAllByStoryID(storyID uint, factionType string) ([]model.Faction, error) {
	var factions []model.Faction
	query := r.db.Where("story_id = ?", storyID)
	if factionType != "" {
		query = query.Where("type = ?", factionType)
	}
	err := query.Order("id ASC").Find(&factions).Error
	if err != nil {
		return nil, err
	}
	return factions, nil
}

func (r *factionRepository) FindByID(id uint) (*model.Faction, error) {
	var faction model.Faction
	err := r.db.First(&faction, id).Error
	if err != nil {
		return nil, err
	}
	return &faction, nil
}

func (r *factionRepository) Update(faction *model.Faction) error {
	return r.db.Save(faction).Error
}

func (r *factionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Faction{}, id).Error
}
