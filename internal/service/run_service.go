package service

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/repository"
	"context"

	"gorm.io/gorm"
)

type RunService interface {
	GetAutomationByID(ctx context.Context, automationID string) (*model.RunAutomation, error)
}

type runService struct {
	automationRepo          repository.AutomationRepository
	automationActionRepo    repository.AutomationActionRepository
	automationConditionRepo repository.AutomationConditionRepository
}

func NewRunService(db *gorm.DB) RunService {
	return &runService{
		automationRepo:          repository.NewAutomationRepository(db),
		automationActionRepo:    repository.NewAutomationActionRepository(db),
		automationConditionRepo: repository.NewAutomationConditionRepository(db),
	}
}

func (s *runService) GetAutomationByID(ctx context.Context, automationID string) (*model.RunAutomation, error) {
	row, err := s.automationRepo.GetByID(ctx, automationID)
	if err != nil {
		return nil, err
	}

	return row, nil
}
