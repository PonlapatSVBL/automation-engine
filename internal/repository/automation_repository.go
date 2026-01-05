package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type AutomationRepository interface {
	GetByID(ctx context.Context, id string) (*model.RunAutomation, error)
}

type automationRepository struct {
	db *gorm.DB
}

func NewAutomationRepository(db *gorm.DB) AutomationRepository {
	return &automationRepository{db: db}
}

func (r *automationRepository) GetByID(ctx context.Context, id string) (*model.RunAutomation, error) {
	q := query.Use(r.db).RunAutomation
	return q.WithContext(ctx).Where(q.AutomationID.Eq(id)).First()
}
