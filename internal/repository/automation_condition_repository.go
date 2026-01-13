package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type AutomationConditionRepository interface {
	ListByGroupIDs(ctx context.Context, groupIDs []string) ([]*model.RunAutomationCondition, error)
}

type automationConditionRepository struct {
	BaseRepository
}

func NewAutomationConditionRepository(db *gorm.DB) AutomationConditionRepository {
	return &automationConditionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *automationConditionRepository) ListByGroupIDs(ctx context.Context, groupIDs []string) ([]*model.RunAutomationCondition, error) {
	q := query.Use(r.Executor(ctx)).RunAutomationCondition
	db := q.WithContext(ctx)

	db = db.Where(q.AutomationConditionGroupID.In(groupIDs...))

	return db.Find()
}
