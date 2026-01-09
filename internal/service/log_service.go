package service

import (
	"automation-engine/internal/repository"
)

type LogService interface{}

type logService struct {
	txManager               repository.TransactionManager
	automationExecutionRepo repository.AutomationExecutionRepository
}

func NewLogService(
	txManager repository.TransactionManager,
	automationExecutionRepo repository.AutomationExecutionRepository,
) LogService {
	return &logService{
		txManager:               txManager,
		automationExecutionRepo: automationExecutionRepo,
	}
}
