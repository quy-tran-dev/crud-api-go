package world

import (
	"strings"

	"crud-api/internal/model"
	"crud-api/internal/repository"
)

type FactionService interface {
	CreateFaction(storyID, userID uint, req model.CreateFactionRequest) (*model.Faction, error)
	GetFactions(storyID, userID uint, factionType string) ([]model.Faction, error)
	GetFactionByID(id, storyID, userID uint) (*model.Faction, error)
	UpdateFaction(id, storyID, userID uint, req model.UpdateFactionRequest) (*model.Faction, error)
	DeleteFaction(id, storyID, userID uint) error
}

type factionService struct {
	factionRepo repository.FactionRepository
	storyRepo   repository.StoryRepository
}

func NewFactionService(factionRepo repository.FactionRepository, storyRepo repository.StoryRepository) FactionService {
	return &factionService{
		factionRepo: factionRepo,
		storyRepo:   storyRepo,
	}
}

func (s *factionService) CreateFaction(storyID, userID uint, req model.CreateFactionRequest) (*model.Faction, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	faction := &model.Faction{
		StoryID:       storyID,
		Name:          strings.TrimSpace(req.Name),
		Type:          strings.TrimSpace(req.Type),
		HeadquarterID: req.HeadquarterID,
		LeaderID:      req.LeaderID,
		PowerTier:     strings.TrimSpace(req.PowerTier),
		Description:   req.Description,
		Diplomacy:     req.Diplomacy,
	}

	if err := s.factionRepo.Create(faction); err != nil {
		return nil, err
	}
	return faction, nil
}

func (s *factionService) GetFactions(storyID, userID uint, factionType string) ([]model.Faction, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}
	return s.factionRepo.FindAllByStoryID(storyID, factionType)
}

func (s *factionService) GetFactionByID(id, storyID, userID uint) (*model.Faction, error) {
	if _, err := s.storyRepo.FindByID(storyID, userID); err != nil {
		return nil, ErrStoryNotFound
	}

	faction, err := s.factionRepo.FindByID(id)
	if err != nil || faction.StoryID != storyID {
		return nil, ErrFactionNotFound
	}
	return faction, nil
}

func (s *factionService) UpdateFaction(id, storyID, userID uint, req model.UpdateFactionRequest) (*model.Faction, error) {
	faction, err := s.GetFactionByID(id, storyID, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		faction.Name = strings.TrimSpace(*req.Name)
	}
	if req.Type != nil {
		faction.Type = strings.TrimSpace(*req.Type)
	}
	if req.HeadquarterID != nil {
		faction.HeadquarterID = req.HeadquarterID
	}
	if req.LeaderID != nil {
		faction.LeaderID = req.LeaderID
	}
	if req.PowerTier != nil {
		faction.PowerTier = strings.TrimSpace(*req.PowerTier)
	}
	if req.Description != nil {
		faction.Description = *req.Description
	}
	if req.Diplomacy != nil {
		faction.Diplomacy = req.Diplomacy
	}

	if err := s.factionRepo.Update(faction); err != nil {
		return nil, err
	}
	return faction, nil
}

func (s *factionService) DeleteFaction(id, storyID, userID uint) error {
	if _, err := s.GetFactionByID(id, storyID, userID); err != nil {
		return err
	}
	return s.factionRepo.Delete(id)
}
