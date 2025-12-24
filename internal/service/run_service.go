package service

import (
	"automation-engine/internal/repository"

	"gorm.io/gorm"
)

type RunService interface{}

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
