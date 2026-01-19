package service

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/repository"
	"context"
	"time"
)

type LogService interface {
	GenerateLogID() string
	Upsert(ctx context.Context, log *model.LogAutomationExecution) error
	DeleteLogsBefore(ctx context.Context, t time.Time) error
}

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

func (s *logService) GenerateLogID() string {
	return s.automationExecutionRepo.GenerateLogID()
}

func (s *logService) Upsert(ctx context.Context, log *model.LogAutomationExecution) error {
	return s.automationExecutionRepo.Upsert(ctx, log)
}

func (s *logService) DeleteLogsBefore(ctx context.Context, t time.Time) error {
	return s.automationExecutionRepo.DeleteBefore(ctx, t)
}
