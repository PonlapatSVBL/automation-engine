package service

import (
	"automation-engine/internal/repository"

	"gorm.io/gorm"
)

type LogService interface{}

type logService struct {
	automationExecutionRepo repository.AutomationExecutionRepository
}

func NewLogService(db *gorm.DB) LogService {
	return &logService{
		automationExecutionRepo: repository.NewAutomationExecutionRepository(db),
	}
}
