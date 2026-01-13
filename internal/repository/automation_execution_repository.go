package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"
	"time"

	"gorm.io/gorm"
)

type AutomationExecutionRepository interface {
	Create(ctx context.Context, log *model.LogAutomationExecution) error
}

type automationExecutionRepository struct {
	BaseRepository
}

func NewAutomationExecutionRepository(db *gorm.DB) AutomationExecutionRepository {
	return &automationExecutionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *automationExecutionRepository) Create(ctx context.Context, log *model.LogAutomationExecution) error {
	log.LogID = r.GenerateSortableID(20)

	if log.TriggeredAt.IsZero() {
		log.TriggeredAt = time.Now()
	}

	q := query.Use(r.Executor(ctx)).LogAutomationExecution

	return q.WithContext(ctx).Create(log)
}
