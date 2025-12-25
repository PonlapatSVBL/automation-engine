package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ConditionActionRepository interface {
	List(ctx context.Context, filter model.PolicyConditionAction) ([]*model.PolicyConditionAction, error)
}

type conditionActionRepository struct {
	db *gorm.DB
}

func NewConditionActionRepository(db *gorm.DB) ConditionActionRepository {
	return &conditionActionRepository{db: db}
}

func (r *conditionActionRepository) List(ctx context.Context, filter model.PolicyConditionAction) ([]*model.PolicyConditionAction, error) {
	q := query.Use(r.db).PolicyConditionAction
	db := q.WithContext(ctx)

	return db.Find()
}
