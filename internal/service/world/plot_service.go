package world

import (
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

type PlotService interface {
	CreatePlot(storyID, userID uint, req model.CreatePlotRequest) (*model.Plot, error)
	GetPlots(storyID, userID uint, plotType string) ([]model.Plot, error)
	GetPlotByID(id, storyID, userID uint) (*model.Plot, error)
	UpdatePlot(id, storyID, userID uint, req model.UpdatePlotRequest) (*model.Plot, error)
	DeletePlot(id, storyID, userID uint) error
}

type plotService struct {
	plotRepo  repository.PlotRepository
	storyRepo repository.StoryRepository
}

func NewPlotService(plotRepo repository.PlotRepository, storyRepo repository.StoryRepository) PlotService {
	return &plotService{
		plotRepo:  plotRepo,
		storyRepo: storyRepo,
	}
}

func (s *plotService) CreatePlot(storyID, userID uint, req model.CreatePlotRequest) (*model.Plot, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	order := 1
	if req.OrderIndex != nil {
		order = *req.OrderIndex
	} else {
		maxOrder, err := s.plotRepo.GetMaxOrderIndex(storyID)
		if err == nil {
			order = maxOrder + 1
		}
	}

	plotType := req.Type
	if plotType == "" {
		plotType = model.PlotTypeMain
	}

	plot := &model.Plot{
		StoryID:                   storyID,
		Title:                     strings.TrimSpace(req.Title),
		Type:                      plotType,
		OrderIndex:                order,
		LocationID:                req.LocationID,
		Summary:                   req.Summary,
		InvolvedFactionIDs:        req.InvolvedFactionIDs,
		ParticipatingCharacterIDs: req.ParticipatingCharacterIDs,
		FeaturedItemIDs:           req.FeaturedItemIDs,
		Consequences:              req.Consequences,
	}

	if err := s.plotRepo.Create(plot); err != nil {
		return nil, err
	}
	return plot, nil
}

func (s *plotService) GetPlots(storyID, userID uint, plotType string) ([]model.Plot, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}
	return s.plotRepo.FindAllByStoryID(storyID, plotType)
}

func (s *plotService) GetPlotByID(id, storyID, userID uint) (*model.Plot, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	plot, err := s.plotRepo.FindByID(id)
	if err != nil || plot.StoryID != storyID {
		return nil, ErrPlotNotFound
	}
	return plot, nil
}

func (s *plotService) UpdatePlot(id, storyID, userID uint, req model.UpdatePlotRequest) (*model.Plot, error) {
	plot, err := s.GetPlotByID(id, storyID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		plot.Title = strings.TrimSpace(*req.Title)
	}
	if req.Type != nil {
		plot.Type = strings.TrimSpace(*req.Type)
	}
	if req.OrderIndex != nil {
		plot.OrderIndex = *req.OrderIndex
	}
	if req.LocationID != nil {
		plot.LocationID = req.LocationID
	}
	if req.Summary != nil {
		plot.Summary = *req.Summary
	}
	if req.InvolvedFactionIDs != nil {
		plot.InvolvedFactionIDs = req.InvolvedFactionIDs
	}
	if req.ParticipatingCharacterIDs != nil {
		plot.ParticipatingCharacterIDs = req.ParticipatingCharacterIDs
	}
	if req.FeaturedItemIDs != nil {
		plot.FeaturedItemIDs = req.FeaturedItemIDs
	}
	if req.Consequences != nil {
		plot.Consequences = *req.Consequences
	}

	if err := s.plotRepo.Update(plot); err != nil {
		return nil, err
	}
	return plot, nil
}

func (s *plotService) DeletePlot(id, storyID, userID uint) error {
	if _, err := s.GetPlotByID(id, storyID, userID); err != nil {
		return err
	}
	return s.plotRepo.Delete(id)
}
