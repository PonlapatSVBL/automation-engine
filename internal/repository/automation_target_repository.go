package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type AutomationTargetRepository interface {
	ListByAutomationID(ctx context.Context, automationID string) ([]*model.RunAutomationTarget, error)
}

type automationTargetRepository struct {
	BaseRepository
}

func NewAutomationTargetRepository(db *gorm.DB) AutomationTargetRepository {
	return &automationTargetRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *automationTargetRepository) ListByAutomationID(ctx context.Context, automationID string) ([]*model.RunAutomationTarget, error) {
	q := query.Use(r.Executor(ctx)).RunAutomationTarget
	db := q.WithContext(ctx)

	db = db.Where(q.AutomationID.Eq(automationID))

	return db.Find()
}
