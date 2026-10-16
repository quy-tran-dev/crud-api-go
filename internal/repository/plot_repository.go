package repository

import (
	"crud-api/internal/model"

	"gorm.io/gorm"
)

type PlotRepository interface {
	Create(plot *model.Plot) error
	FindAllByStoryID(storyID uint, plotType string) ([]model.Plot, error)
	FindByID(id uint) (*model.Plot, error)
	Update(plot *model.Plot) error
	Delete(id uint) error
	GetMaxOrderIndex(storyID uint) (int, error)
}

type plotRepository struct {
	db *gorm.DB
}

func NewPlotRepository(db *gorm.DB) PlotRepository {
	return &plotRepository{db: db}
}

func (r *plotRepository) Create(plot *model.Plot) error {
	return r.db.Create(plot).Error
}

func (r *plotRepository) FindAllByStoryID(storyID uint, plotType string) ([]model.Plot, error) {
	var plots []model.Plot
	query := r.db.Where("story_id = ?", storyID)
	if plotType != "" {
		query = query.Where("type = ?", plotType)
	}
	err := query.Order("order_index ASC, id ASC").Find(&plots).Error
	if err != nil {
		return nil, err
	}
	return plots, nil
}

func (r *plotRepository) FindByID(id uint) (*model.Plot, error) {
	var plot model.Plot
	err := r.db.First(&plot, id).Error
	if err != nil {
		return nil, err
	}
	return &plot, nil
}

func (r *plotRepository) Update(plot *model.Plot) error {
	return r.db.Save(plot).Error
}

func (r *plotRepository) Delete(id uint) error {
	return r.db.Delete(&model.Plot{}, id).Error
}

func (r *plotRepository) GetMaxOrderIndex(storyID uint) (int, error) {
	var maxOrder int
	row := r.db.Model(&model.Plot{}).Where("story_id = ?", storyID).Select("COALESCE(MAX(order_index), 0)").Row()
	err := row.Scan(&maxOrder)
	return maxOrder, err
}
