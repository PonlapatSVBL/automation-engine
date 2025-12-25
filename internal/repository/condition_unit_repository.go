package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ConditionUnitRepository interface {
	List(ctx context.Context, filter model.PolicyConditionUnit) ([]*model.PolicyConditionUnit, error)
}

type conditionUnitRepository struct {
	db *gorm.DB
}

func NewConditionUnitRepository(db *gorm.DB) ConditionUnitRepository {
	return &conditionUnitRepository{db: db}
}

func (r *conditionUnitRepository) List(ctx context.Context, filter model.PolicyConditionUnit) ([]*model.PolicyConditionUnit, error) {
	q := query.Use(r.db).PolicyConditionUnit
	db := q.WithContext(ctx)

	return db.Find()
}
