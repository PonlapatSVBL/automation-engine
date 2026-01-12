package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ConditionUnitRepository interface {
	List(ctx context.Context, filter model.PolicyConditionUnit) ([]*model.PolicyConditionUnit, error)
	DeleteByConditionID(ctx context.Context, conditionID string) error
	BulkCreate(ctx context.Context, ops []*model.PolicyConditionUnit) error
	WithTransaction(ctx context.Context, fn func(txRepo ConditionUnitRepository) error) error
}

type conditionUnitRepository struct {
	BaseRepository
}

func NewConditionUnitRepository(db *gorm.DB) ConditionUnitRepository {
	return &conditionUnitRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *conditionUnitRepository) List(ctx context.Context, filter model.PolicyConditionUnit) ([]*model.PolicyConditionUnit, error) {
	q := query.Use(r.Executor(ctx)).PolicyConditionUnit
	db := q.WithContext(ctx)

	return db.Find()
}

func (r *conditionUnitRepository) DeleteByConditionID(ctx context.Context, conditionID string) error {
	return r.Executor(ctx).
		Where("condition_id = ?", conditionID).
		Delete(&model.PolicyConditionUnit{}).Error
}

func (r *conditionUnitRepository) BulkCreate(ctx context.Context, ops []*model.PolicyConditionUnit) error {
	return r.Executor(ctx).
		Create(&ops).Error
}

func (r *conditionUnitRepository) WithTransaction(ctx context.Context, fn func(txRepo ConditionUnitRepository) error) error {
	return r.Executor(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewConditionUnitRepository(tx)
		return fn(txRepo)
	})
}
