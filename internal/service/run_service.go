package service

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/repository"
	"context"
)

type RunService interface {
	GetAutomationByID(ctx context.Context, automationID string) (*model.RunAutomation, error)
	UpdateAutomationByID(ctx context.Context, automation *model.RunAutomation) error
	FetchAndLockTasks(ctx context.Context, limit int) ([]*model.RunAutomation, error)
	MarkTasksCompleted(ctx context.Context, taskIDs []string) error
	BulkUpdateNextRun(ctx context.Context, tasks []*model.RunAutomation) error
}

type runService struct {
	txManager               repository.TransactionManager
	automationRepo          repository.AutomationRepository
	automationActionRepo    repository.AutomationActionRepository
	automationConditionRepo repository.AutomationConditionRepository
}

func NewRunService(
	txManager repository.TransactionManager,
	automationRepo repository.AutomationRepository,
	automationActionRepo repository.AutomationActionRepository,
	automationConditionRepo repository.AutomationConditionRepository,
) RunService {
	return &runService{
		txManager:               txManager,
		automationRepo:          automationRepo,
		automationActionRepo:    automationActionRepo,
		automationConditionRepo: automationConditionRepo,
	}
}

func (s *runService) GetAutomationByID(ctx context.Context, automationID string) (*model.RunAutomation, error) {
	row, err := s.automationRepo.GetByID(ctx, automationID)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (s *runService) UpdateAutomationByID(ctx context.Context, automation *model.RunAutomation) error {
	err := s.automationRepo.Update(ctx, automation)
	if err != nil {
		return err
	}

	return nil
}

func (s *runService) FetchAndLockTasks(ctx context.Context, limit int) ([]*model.RunAutomation, error) {
	var tasks []*model.RunAutomation

	err := s.automationRepo.WithTransaction(ctx, func(txRepo repository.AutomationRepository) error {
		// 1. ดึงงานและ Lock แถวไว้
		lockedTasks, err := txRepo.FetchAndLock(ctx, limit)
		if err != nil {
			return err
		}

		if len(lockedTasks) == 0 {
			return nil
		}

		// 2. เก็บ IDs เพื่อไปอัปเดตสถานะ
		var ids []string
		for _, t := range lockedTasks {
			ids = append(ids, t.AutomationID)
		}

		// 3. เปลี่ยนสถานะเป็น PROCESSING ทันที (จองงาน)
		if err := txRepo.UpdateStatusBatch(ctx, ids, "PROCESSING"); err != nil {
			return err
		}

		tasks = lockedTasks
		return nil
	})

	return tasks, err
}

func (s *runService) MarkTasksCompleted(ctx context.Context, taskIDs []string) error {
	if len(taskIDs) == 0 {
		return nil
	}

	if err := s.automationRepo.UpdateStatusBatch(ctx, taskIDs, "COMPLETE"); err != nil {
		return err
	}

	return nil
}

func (s *runService) BulkUpdateNextRun(ctx context.Context, tasks []*model.RunAutomation) error {
	if len(tasks) == 0 {
		return nil
	}

	return s.automationRepo.BulkUpdateNextRun(ctx, tasks)
}
