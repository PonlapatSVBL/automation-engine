package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type OperatorRepository interface {
	List(ctx context.Context, filter model.DefOperator) ([]*model.DefOperator, error)
}

type operatorRepository struct {
	BaseRepository
}

func NewOperatorRepository(db *gorm.DB) OperatorRepository {
	return &operatorRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *operatorRepository) List(ctx context.Context, filter model.DefOperator) ([]*model.DefOperator, error) {
	q := query.Use(r.Executor(ctx)).DefOperator
	db := q.WithContext(ctx)

	return db.Find()
}
