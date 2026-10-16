package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type CharacterRepository interface {
	Create(character *model.Character) error
	FindAllByStoryID(storyID uint, factionID *uint, role, level string) ([]model.Character, error)
	FindByID(id uint) (*model.Character, error)
	FindByIDs(storyID uint, ids []uint) ([]model.Character, error)
	Update(character *model.Character) error
	Delete(id uint) error
}

type characterRepository struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) CharacterRepository {
	return &characterRepository{db: db}
}

func (r *characterRepository) Create(character *model.Character) error {
	return r.db.Create(character).Error
}

func (r *characterRepository) FindAllByStoryID(storyID uint, factionID *uint, role, level string) ([]model.Character, error) {
	var characters []model.Character
	query := r.db.Where("story_id = ?", storyID)

	if factionID != nil {
		query = query.Where("faction_id = ?", *factionID)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if level != "" {
		query = query.Where("current_level = ?", level)
	}

	err := query.Order("id ASC").Find(&characters).Error
	if err != nil {
		return nil, err
	}
	return characters, nil
}

func (r *characterRepository) FindByID(id uint) (*model.Character, error) {
	var character model.Character
	err := r.db.First(&character, id).Error
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func (r *characterRepository) FindByIDs(storyID uint, ids []uint) ([]model.Character, error) {
	if len(ids) == 0 {
		return []model.Character{}, nil
	}
	var characters []model.Character
	err := r.db.Where("story_id = ? AND id IN ?", storyID, ids).Find(&characters).Error
	if err != nil {
		return nil, err
	}
	return characters, nil
}

func (r *characterRepository) Update(character *model.Character) error {
	return r.db.Save(character).Error
}

func (r *characterRepository) Delete(id uint) error {
	return r.db.Delete(&model.Character{}, id).Error
}
