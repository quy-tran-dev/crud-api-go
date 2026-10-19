package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *model.SystemCategory) error
	FindCategoryByID(id uint) (*model.SystemCategory, error)
	FindCategoriesByStoryID(storyID uint) ([]model.SystemCategory, error)
	FindCategoryByCode(storyID uint, code string) (*model.SystemCategory, error)
	UpdateCategory(category *model.SystemCategory) error
	DeleteCategory(id uint) error

	CreateCategoryItem(item *model.SystemCategoryItem) error
	FindCategoryItemByID(id uint) (*model.SystemCategoryItem, error)
	FindCategoryItems(categoryID uint) ([]model.SystemCategoryItem, error)
	FindCategoryItemByName(categoryID uint, name string) (*model.SystemCategoryItem, error)
	UpdateCategoryItem(item *model.SystemCategoryItem) error
	DeleteCategoryItem(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *model.SystemCategory) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindCategoriesByStoryID(storyID uint) ([]model.SystemCategory, error) {
	var categories []model.SystemCategory
	err := r.db.Where("story_id = ?", storyID).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC, id ASC")
		}).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) FindCategoryByCode(storyID uint, code string) (*model.SystemCategory, error) {
	var category model.SystemCategory
	err := r.db.Where("story_id = ? AND code = ?", storyID, code).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) CreateCategoryItem(item *model.SystemCategoryItem) error {
	return r.db.Create(item).Error
}

func (r *categoryRepository) FindCategoryItems(categoryID uint) ([]model.SystemCategoryItem, error) {
	var items []model.SystemCategoryItem
	err := r.db.Where("category_id = ?", categoryID).Order("order_index ASC, id ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *categoryRepository) FindCategoryItemByName(categoryID uint, name string) (*model.SystemCategoryItem, error) {
	var item model.SystemCategoryItem
	err := r.db.Where("category_id = ? AND LOWER(name) = LOWER(?)", categoryID, name).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *categoryRepository) FindCategoryByID(id uint) (*model.SystemCategory, error) {
	var category model.SystemCategory
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_index ASC, id ASC")
	}).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) UpdateCategory(category *model.SystemCategory) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) FindCategoryItemByID(id uint) (*model.SystemCategoryItem, error) {
	var item model.SystemCategoryItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *categoryRepository) UpdateCategoryItem(item *model.SystemCategoryItem) error {
	return r.db.Save(item).Error
}

func (r *categoryRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&model.SystemCategory{}, id).Error
}

func (r *categoryRepository) DeleteCategoryItem(id uint) error {
	return r.db.Delete(&model.SystemCategoryItem{}, id).Error
}
