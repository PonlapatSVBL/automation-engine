package service

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/dto"
	"automation-engine/internal/repository"
	"context"
	"time"
)

type RunService interface {
	GetAutomationByID(ctx context.Context, automationID string) (*model.RunAutomation, error)
	GetAutomationSnapshot(ctx context.Context, automationID string) (*dto.AutomationSnapshot, error)
	UpdateAutomationByID(ctx context.Context, automation *model.RunAutomation) error
	FetchAndLockTasks(ctx context.Context, runTime time.Time, limit int) ([]*model.RunAutomation, error)
	MarkTasksCompleted(ctx context.Context, taskIDs []string) error
	BulkUpdateNextRun(ctx context.Context, tasks []*model.RunAutomation) error
}

type runService struct {
	txManager                    repository.TransactionManager
	automationRepo               repository.AutomationRepository
	automationActionRepo         repository.AutomationActionRepository
	automationConditionGroupRepo repository.AutomationConditionGroupRepository
	automationConditionRepo      repository.AutomationConditionRepository
	automationTargetRepo         repository.AutomationTargetRepository
	automationExecutionRepo      repository.AutomationExecutionRepository
}

func NewRunService(
	txManager repository.TransactionManager,
	automationRepo repository.AutomationRepository,
	automationActionRepo repository.AutomationActionRepository,
	automationConditionGroupRepo repository.AutomationConditionGroupRepository,
	automationConditionRepo repository.AutomationConditionRepository,
	automationTargetRepo repository.AutomationTargetRepository,
	automationExecutionRepo repository.AutomationExecutionRepository,
) RunService {
	return &runService{
		txManager:                    txManager,
		automationRepo:               automationRepo,
		automationActionRepo:         automationActionRepo,
		automationConditionGroupRepo: automationConditionGroupRepo,
		automationConditionRepo:      automationConditionRepo,
		automationTargetRepo:         automationTargetRepo,
		automationExecutionRepo:      automationExecutionRepo,
	}
}

func (s *runService) GetAutomationByID(ctx context.Context, automationID string) (*model.RunAutomation, error) {
	row, err := s.automationRepo.GetByID(ctx, automationID)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (s *runService) GetAutomationSnapshot(ctx context.Context, automationID string) (*dto.AutomationSnapshot, error) {
	snapshot := &dto.AutomationSnapshot{}

	automation, err := s.automationRepo.GetByID(ctx, automationID)
	if err != nil {
		return nil, err
	}

	conditionGroups, err := s.automationConditionGroupRepo.ListByAutomationID(ctx, automationID)
	if err != nil {
		return nil, err
	}

	var groupIDs []string
	for _, group := range conditionGroups {
		groupIDs = append(groupIDs, group.AutomationConditionGroupID)
	}

	conditions, err := s.automationConditionRepo.ListByGroupIDs(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	actions, err := s.automationActionRepo.ListByAutomationID(ctx, automationID)
	if err != nil {
		return nil, err
	}

	targets, err := s.automationTargetRepo.ListByAutomationID(ctx, automationID)
	if err != nil {
		return nil, err
	}

	snapshot.Automation = automation
	snapshot.ConditionGroups = conditionGroups
	snapshot.Conditions = conditions
	snapshot.Actions = actions
	snapshot.Targets = targets

	return snapshot, nil
}

func (s *runService) UpdateAutomationByID(ctx context.Context, automation *model.RunAutomation) error {
	err := s.automationRepo.Update(ctx, automation)
	if err != nil {
		return err
	}

	return nil
}

func (s *runService) FetchAndLockTasks(ctx context.Context, runTime time.Time, limit int) ([]*model.RunAutomation, error) {
	var tasks []*model.RunAutomation

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. ดึงงานและ Lock แถวไว้
		lockedTasks, err := s.automationRepo.FetchAndLock(txCtx, runTime, limit)
		if err != nil {
			return err
		}

		if len(lockedTasks) == 0 {
			return nil
		}

		// 2. เก็บ IDs เพื่อไปอัพเดทสถานะ
		var ids []string
		for _, t := range lockedTasks {
			ids = append(ids, t.AutomationID)
		}

		// 3. เปลี่ยนสถานะเป็น LOCKED ทันที (จองงาน)
		if err := s.automationRepo.UpdateStatusBatch(txCtx, ids, "LOCKED"); err != nil {
			return err
		}

		tasks = lockedTasks
		return nil
	})
	/* err := s.automationRepo.WithTransaction(ctx, func(txRepo repository.AutomationRepository) error {
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
	}) */

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

	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return s.automationRepo.BulkUpdateNextRun(txCtx, tasks)
	})
}
