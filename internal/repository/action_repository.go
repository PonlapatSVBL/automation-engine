package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"

	"gorm.io/gorm"
)

type ActionRepository interface {
	GetByID(ctx context.Context, id string) (*model.DefAction, error)
	Create(ctx context.Context, action *model.DefAction) error
	Update(ctx context.Context, action *model.DefAction) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter model.DefAction) ([]*model.DefAction, error)
	ListByActionIDs(ctx context.Context, actionIDs []string) ([]*model.DefAction, error)
}

type actionRepository struct {
	BaseRepository
}

func NewActionRepository(db *gorm.DB) ActionRepository {
	return &actionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *actionRepository) GetByID(ctx context.Context, id string) (*model.DefAction, error) {
	q := query.Use(r.Executor(ctx)).DefAction
	return q.WithContext(ctx).Where(q.ActionID.Eq(id)).First()
}

func (r *actionRepository) Create(ctx context.Context, action *model.DefAction) error {
	q := query.Use(r.Executor(ctx)).DefAction
	return q.WithContext(ctx).Create(action)
}

func (r *actionRepository) Update(ctx context.Context, action *model.DefAction) error {
	q := query.Use(r.Executor(ctx)).DefAction
	// ใช้ Select("*") เพื่อบังคับให้อัปเดตทุกฟิลด์รวมถึงค่าว่าง หรือระบุฟิลด์ที่ต้องการ
	_, err := q.WithContext(ctx).Where(q.ActionID.Eq(action.ActionID)).Updates(action)
	return err
}

func (r *actionRepository) Delete(ctx context.Context, id string) error {
	q := query.Use(r.Executor(ctx)).DefAction
	// GORM Gen จะจัดการ Soft Delete ให้โดยอัตโนมัติหากใน Model มีฟิลด์ DeletedAt
	_, err := q.WithContext(ctx).Where(q.ActionID.Eq(id)).Delete()
	return err
}

func (r *actionRepository) List(ctx context.Context, filter model.DefAction) ([]*model.DefAction, error) {
	q := query.Use(r.Executor(ctx)).DefAction
	db := q.WithContext(ctx)

	// Dynamic Filtering
	if filter.ActionCode != "" {
		db = db.Where(q.ActionCode.Eq(filter.ActionCode))
	}
	if filter.ActionType != "" {
		db = db.Where(q.ActionType.Eq(filter.ActionType))
	}
	if filter.Status != "" {
		db = db.Where(q.Status.Eq(filter.Status))
	}

	return db.Find()
}

func (r *actionRepository) ListByActionIDs(ctx context.Context, actionIDs []string) ([]*model.DefAction, error) {
	q := query.Use(r.Executor(ctx)).DefAction
	db := q.WithContext(ctx)

	db = db.Where(q.ActionID.In(actionIDs...))

	return db.Find()
}
