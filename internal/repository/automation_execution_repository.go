package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AutomationExecutionRepository interface {
	GenerateLogID() string
	Create(ctx context.Context, log *model.LogAutomationExecution) error
	Upsert(ctx context.Context, log *model.LogAutomationExecution) error
}

type automationExecutionRepository struct {
	BaseRepository
}

func NewAutomationExecutionRepository(db *gorm.DB) AutomationExecutionRepository {
	return &automationExecutionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *automationExecutionRepository) GenerateLogID() string {
	return r.GenerateSortableID(20)
}

func (r *automationExecutionRepository) Create(ctx context.Context, log *model.LogAutomationExecution) error {
	log.LogID = r.GenerateSortableID(20)

	if log.TriggeredAt.IsZero() {
		log.TriggeredAt = time.Now()
	}

	q := query.Use(r.Executor(ctx)).LogAutomationExecution

	return q.WithContext(ctx).Create(log)
}

func (r *automationExecutionRepository) Upsert(ctx context.Context, log *model.LogAutomationExecution) error {
	return r.Executor(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "log_id"}},
		UpdateAll: true,
	}).Create(log).Error
}
