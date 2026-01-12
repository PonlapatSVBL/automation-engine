package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ConditionRepository interface {
	List(ctx context.Context, filter model.DefCondition) ([]*model.DefCondition, error)
}

type conditionRepository struct {
	BaseRepository
}

func NewConditionRepository(db *gorm.DB) ConditionRepository {
	return &conditionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *conditionRepository) List(ctx context.Context, filter model.DefCondition) ([]*model.DefCondition, error) {
	q := query.Use(r.Executor(ctx)).DefCondition
	db := q.WithContext(ctx)

	// Dynamic Filtering
	if filter.ConditionCode != "" {
		db = db.Where(q.ConditionCode.Eq(filter.ConditionCode))
	}
	if filter.ConditionName != "" {
		db = db.Where(q.ConditionName.Eq(filter.ConditionName))
	}
	if filter.Status != "" {
		db = db.Where(q.Status.Eq(filter.Status))
	}

	return db.Find()
}
