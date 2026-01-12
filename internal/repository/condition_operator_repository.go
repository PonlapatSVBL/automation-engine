package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ConditionOperatorRepository interface {
	List(ctx context.Context, filter model.PolicyConditionOperator) ([]*model.PolicyConditionOperator, error)
	DeleteByConditionID(ctx context.Context, conditionID string) error
	BulkCreate(ctx context.Context, ops []*model.PolicyConditionOperator) error
	// WithTransaction(ctx context.Context, fn func(txRepo ConditionOperatorRepository) error) error
}

type conditionOperatorRepository struct {
	BaseRepository
}

func NewConditionOperatorRepository(db *gorm.DB) ConditionOperatorRepository {
	return &conditionOperatorRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *conditionOperatorRepository) List(ctx context.Context, filter model.PolicyConditionOperator) ([]*model.PolicyConditionOperator, error) {
	q := query.Use(r.Executor(ctx)).PolicyConditionOperator
	db := q.WithContext(ctx)

	return db.Find()
}

func (r *conditionOperatorRepository) DeleteByConditionID(ctx context.Context, conditionID string) error {
	return r.Executor(ctx).
		Where("condition_id = ?", conditionID).
		Delete(&model.PolicyConditionOperator{}).Error
}

func (r *conditionOperatorRepository) BulkCreate(ctx context.Context, ops []*model.PolicyConditionOperator) error {
	return r.Executor(ctx).
		Create(&ops).Error
}

/* func (r *conditionOperatorRepository) WithTransaction(ctx context.Context, fn func(txRepo ConditionOperatorRepository) error) error {
	return r.executor(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewConditionOperatorRepository(tx)
		return fn(txRepo)
	})
} */
