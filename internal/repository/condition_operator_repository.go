package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ConditionOperatorRepository interface {
	List(ctx context.Context, filter model.PolicyConditionOperator) ([]*model.PolicyConditionOperator, error)
}

type conditionOperatorRepository struct {
	db *gorm.DB
}

func NewConditionOperatorRepository(db *gorm.DB) ConditionOperatorRepository {
	return &conditionOperatorRepository{db: db}
}

func (r *conditionOperatorRepository) List(ctx context.Context, filter model.PolicyConditionOperator) ([]*model.PolicyConditionOperator, error) {
	q := query.Use(r.db).PolicyConditionOperator
	db := q.WithContext(ctx)

	return db.Find()
}
