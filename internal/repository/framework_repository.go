package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type FrameworkRepository interface {
	Create(framework *model.FrameworkTemplate) error
	FindByID(id uint) (*model.FrameworkTemplate, error)
	FindByUserID(userID uint) ([]model.FrameworkTemplate, error)
	FindPublic(genre string, search string) ([]model.FrameworkTemplate, error)
	Update(framework *model.FrameworkTemplate) error
	Delete(id uint) error
	IncrementUsage(id uint) error
}

type frameworkRepository struct {
	db *gorm.DB
}

func NewFrameworkRepository(db *gorm.DB) FrameworkRepository {
	return &frameworkRepository{db: db}
}

func (r *frameworkRepository) Create(framework *model.FrameworkTemplate) error {
	return r.db.Create(framework).Error
}

func (r *frameworkRepository) FindByID(id uint) (*model.FrameworkTemplate, error) {
	var framework model.FrameworkTemplate
	err := r.db.Preload("Author").First(&framework, id).Error
	if err != nil {
		return nil, err
	}
	return &framework, nil
}

func (r *frameworkRepository) FindByUserID(userID uint) ([]model.FrameworkTemplate, error) {
	var list []model.FrameworkTemplate
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *frameworkRepository) FindPublic(genre string, search string) ([]model.FrameworkTemplate, error) {
	var list []model.FrameworkTemplate
	query := r.db.Preload("Author").Where("is_public = ?", true)

	if genre != "" && genre != "ALL" {
		query = query.Where("LOWER(genre) = LOWER(?)", genre)
	}

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR genre ILIKE ?", s, s, s)
	}

	err := query.Order("usage_count DESC, created_at DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *frameworkRepository) Update(framework *model.FrameworkTemplate) error {
	return r.db.Save(framework).Error
}

func (r *frameworkRepository) Delete(id uint) error {
	return r.db.Delete(&model.FrameworkTemplate{}, id).Error
}

func (r *frameworkRepository) IncrementUsage(id uint) error {
	return r.db.Model(&model.FrameworkTemplate{}).Where("id = ?", id).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error
}
