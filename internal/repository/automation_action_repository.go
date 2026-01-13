package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type AutomationActionRepository interface {
	ListByAutomationID(ctx context.Context, automationID string) ([]*model.RunAutomationAction, error)
}

type automationActionRepository struct {
	BaseRepository
}

func NewAutomationActionRepository(db *gorm.DB) AutomationActionRepository {
	return &automationActionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *automationActionRepository) ListByAutomationID(ctx context.Context, automationID string) ([]*model.RunAutomationAction, error) {
	q := query.Use(r.Executor(ctx)).RunAutomationAction
	db := q.WithContext(ctx)

	db = db.Where(q.AutomationID.Eq(automationID))

	return db.Find()
}
