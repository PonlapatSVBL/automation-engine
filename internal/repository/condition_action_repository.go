package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ConditionActionRepository interface {
	List(ctx context.Context, filter model.PolicyConditionAction) ([]*model.PolicyConditionAction, error)
	DeleteByConditionID(ctx context.Context, conditionID string) error
	BulkCreate(ctx context.Context, ops []*model.PolicyConditionAction) error
	WithTransaction(ctx context.Context, fn func(txRepo ConditionActionRepository) error) error
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

func (r *conditionActionRepository) DeleteByConditionID(ctx context.Context, conditionID string) error {
	return r.db.WithContext(ctx).
		Where("condition_id = ?", conditionID).
		Delete(&model.PolicyConditionAction{}).Error
}

func (r *conditionActionRepository) BulkCreate(ctx context.Context, ops []*model.PolicyConditionAction) error {
	return r.db.WithContext(ctx).
		Create(&ops).Error
}

func (r *conditionActionRepository) WithTransaction(ctx context.Context, fn func(txRepo ConditionActionRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewConditionActionRepository(tx)
		return fn(txRepo)
	})
}
