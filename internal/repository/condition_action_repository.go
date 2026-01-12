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
}

type conditionActionRepository struct {
	BaseRepository
}

func NewConditionActionRepository(db *gorm.DB) ConditionActionRepository {
	return &conditionActionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *conditionActionRepository) List(ctx context.Context, filter model.PolicyConditionAction) ([]*model.PolicyConditionAction, error) {
	q := query.Use(r.Executor(ctx)).PolicyConditionAction
	db := q.WithContext(ctx)

	return db.Find()
}

func (r *conditionActionRepository) DeleteByConditionID(ctx context.Context, conditionID string) error {
	return r.Executor(ctx).
		Where("condition_id = ?", conditionID).
		Delete(&model.PolicyConditionAction{}).Error
}

func (r *conditionActionRepository) BulkCreate(ctx context.Context, ops []*model.PolicyConditionAction) error {
	return r.Executor(ctx).
		Create(&ops).Error
}
