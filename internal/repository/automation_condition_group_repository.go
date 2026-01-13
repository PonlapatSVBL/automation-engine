package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type AutomationConditionGroupRepository interface {
	ListByAutomationID(ctx context.Context, automationID string) ([]*model.RunAutomationConditionGroup, error)
}

type automationConditionGroupRepository struct {
	BaseRepository
}

func NewAutomationConditionGroupRepository(db *gorm.DB) AutomationConditionGroupRepository {
	return &automationConditionGroupRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *automationConditionGroupRepository) ListByAutomationID(ctx context.Context, automationID string) ([]*model.RunAutomationConditionGroup, error) {
	q := query.Use(r.Executor(ctx)).RunAutomationConditionGroup
	db := q.WithContext(ctx)

	db = db.Where(q.AutomationID.Eq(automationID))

	return db.Find()
}
