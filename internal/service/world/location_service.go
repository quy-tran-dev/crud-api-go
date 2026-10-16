package world

import (
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

type LocationService interface {
	CreateLocation(storyID, userID uint, req model.CreateLocationRequest) (*model.Location, error)
	GetLocations(storyID, userID uint, treeMode bool) ([]model.Location, error)
	GetLocationByID(id, storyID, userID uint) (*model.Location, error)
	UpdateLocation(id, storyID, userID uint, req model.UpdateLocationRequest) (*model.Location, error)
	DeleteLocation(id, storyID, userID uint) error
}

type locationService struct {
	locationRepo repository.LocationRepository
	storyRepo    repository.StoryRepository
}

func NewLocationService(locationRepo repository.LocationRepository, storyRepo repository.StoryRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
		storyRepo:    storyRepo,
	}
}

func (s *locationService) CreateLocation(storyID, userID uint, req model.CreateLocationRequest) (*model.Location, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	location := &model.Location{
		StoryID:     storyID,
		ParentID:    req.ParentID,
		Name:        strings.TrimSpace(req.Name),
		Type:        strings.TrimSpace(req.Type),
		DangerLevel: req.DangerLevel,
		Description: req.Description,
	}

	if err := s.locationRepo.Create(location); err != nil {
		return nil, err
	}
	return location, nil
}

func (s *locationService) GetLocations(storyID, userID uint, treeMode bool) ([]model.Location, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	if treeMode {
		return s.locationRepo.FindTreeByStoryID(storyID)
	}
	return s.locationRepo.FindAllByStoryID(storyID)
}

func (s *locationService) GetLocationByID(id, storyID, userID uint) (*model.Location, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	location, err := s.locationRepo.FindByID(id)
	if err != nil || location.StoryID != storyID {
		return nil, ErrLocationNotFound
	}
	return location, nil
}

func (s *locationService) UpdateLocation(id, storyID, userID uint, req model.UpdateLocationRequest) (*model.Location, error) {
	location, err := s.GetLocationByID(id, storyID, userID)
	if err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		// Tránh tình trạng gán cha chính là bản thân
		if *req.ParentID == id {
			req.ParentID = nil
		}
		location.ParentID = req.ParentID
	}
	if req.Name != nil {
		location.Name = strings.TrimSpace(*req.Name)
	}
	if req.Type != nil {
		location.Type = strings.TrimSpace(*req.Type)
	}
	if req.DangerLevel != nil {
		location.DangerLevel = *req.DangerLevel
	}
	if req.Description != nil {
		location.Description = *req.Description
	}

	if err := s.locationRepo.Update(location); err != nil {
		return nil, err
	}
	return location, nil
}

func (s *locationService) DeleteLocation(id, storyID, userID uint) error {
	if _, err := s.GetLocationByID(id, storyID, userID); err != nil {
		return err
	}
	return s.locationRepo.Delete(id)
}
