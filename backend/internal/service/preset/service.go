package preset

import (
	"context"
	"errors"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

type Service struct {
	repo *repository.PresetRepository
}

func NewService(repo *repository.PresetRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePreset(ctx context.Context, userID string, req model.Preset) (*model.Preset, error) {
	if req.Name == "" {
		return nil, errors.New("preset name is required")
	}
	req.CreatorID = userID
	if req.Blocks == nil {
		req.Blocks = []model.PresetBlock{}
	}
	if req.GenParams == nil {
		req.GenParams = map[string]interface{}{}
	}
	if err := s.repo.Create(ctx, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *Service) UpdatePreset(ctx context.Context, userID, presetID string, req model.Preset) (*model.Preset, error) {
	existing, err := s.repo.FindByID(ctx, presetID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("preset not found")
	}
	if existing.CreatorID != userID {
		return nil, errors.New("forbidden")
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.ModelKey = req.ModelKey
	if req.Blocks != nil {
		existing.Blocks = req.Blocks
	} else {
		existing.Blocks = []model.PresetBlock{}
	}
	if req.GenParams != nil {
		existing.GenParams = req.GenParams
	} else {
		existing.GenParams = map[string]interface{}{}
	}
	existing.IsPublic = req.IsPublic

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeletePreset(ctx context.Context, userID, presetID string) error {
	existing, err := s.repo.FindByID(ctx, presetID)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil // idempotent
	}
	if existing.CreatorID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Delete(ctx, presetID)
}

func (s *Service) GetPreset(ctx context.Context, userID, presetID string) (*model.Preset, error) {
	preset, err := s.repo.FindByID(ctx, presetID)
	if err != nil {
		return nil, err
	}
	if preset == nil {
		return nil, errors.New("preset not found")
	}
	// Allow if public or owner
	if !preset.IsPublic && preset.CreatorID != userID {
		return nil, errors.New("forbidden")
	}
	return preset, nil
}

func (s *Service) ListUserPresets(ctx context.Context, userID string) ([]model.Preset, error) {
	return s.repo.ListByCreator(ctx, userID)
}

func (s *Service) PublishPreset(ctx context.Context, userID, presetID string, isPublic bool) (*model.Preset, error) {
	existing, err := s.repo.FindByID(ctx, presetID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("preset not found")
	}
	if existing.CreatorID != userID {
		return nil, errors.New("forbidden")
	}
	if err := s.repo.UpdatePublic(ctx, presetID, isPublic); err != nil {
		return nil, err
	}
	existing.IsPublic = isPublic
	return existing, nil
}

func (s *Service) ListPublicPresets(ctx context.Context, limit int) ([]model.Preset, error) {
	return s.repo.ListPublic(ctx, limit)
}
