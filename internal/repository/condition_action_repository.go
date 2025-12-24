package repository

import (
	"gorm.io/gorm"
)

type ConditionActionRepository interface {
	// GetByID(ctx context.Context, id string) (*model.PolicyConditionAction, error)
	// Create(ctx context.Context, action *model.PolicyConditionAction) error
	// Update(ctx context.Context, action *model.PolicyConditionAction) error
	// Delete(ctx context.Context, id string) error
	// List(ctx context.Context, filter model.PolicyConditionAction) ([]*model.PolicyConditionAction, error)
}

type conditionActionRepository struct {
	db *gorm.DB
}

func NewConditionActionRepository(db *gorm.DB) ConditionActionRepository {
	return &conditionActionRepository{db: db}
}
