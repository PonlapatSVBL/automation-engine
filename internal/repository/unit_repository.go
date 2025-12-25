package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type UnitRepository interface {
	List(ctx context.Context, filter model.DefUnit) ([]*model.DefUnit, error)
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) List(ctx context.Context, filter model.DefUnit) ([]*model.DefUnit, error) {
	q := query.Use(r.db).DefUnit
	db := q.WithContext(ctx)

	return db.Find()
}
