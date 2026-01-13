package repository

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/domain/query"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AutomationRepository interface {
	GetByID(ctx context.Context, id string) (*model.RunAutomation, error)
	Update(ctx context.Context, action *model.RunAutomation) error
	FetchAndLock(ctx context.Context, runTime time.Time, limit int) ([]*model.RunAutomation, error)
	UpdateStatusBatch(ctx context.Context, ids []string, status string) error
	BulkUpdateNextRun(ctx context.Context, tasks []*model.RunAutomation) error
}

type automationRepository struct {
	BaseRepository
}

func NewAutomationRepository(db *gorm.DB) AutomationRepository {
	return &automationRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *automationRepository) GetByID(ctx context.Context, id string) (*model.RunAutomation, error) {
	q := query.Use(r.Executor(ctx)).RunAutomation
	return q.WithContext(ctx).Where(q.AutomationID.Eq(id)).First()
}

func (r *automationRepository) Update(ctx context.Context, action *model.RunAutomation) error {
	q := query.Use(r.Executor(ctx)).RunAutomation
	_, err := q.WithContext(ctx).Where(q.AutomationID.Eq(action.AutomationID)).Updates(action)
	return err
}

func (r *automationRepository) FetchAndLock(ctx context.Context, runTime time.Time, limit int) ([]*model.RunAutomation, error) {
	var results []*model.RunAutomation
	q := query.Use(r.Executor(ctx)).RunAutomation

	err := r.Executor(ctx).WithContext(ctx).
		Model(&model.RunAutomation{}).
		Where(q.NextRunTime.Lte(runTime)).
		Where(q.Status.Eq("PENDING")).
		Limit(limit).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Find(&results).Error

	return results, err
}

func (r *automationRepository) UpdateStatusBatch(ctx context.Context, ids []string, status string) error {
	q := query.Use(r.Executor(ctx)).RunAutomation

	_, err := q.WithContext(ctx).
		Where(q.AutomationID.In(ids...)).
		Updates(&model.RunAutomation{
			Status:  status,
			LastUpd: time.Now(),
		})
	return err
}

func (r *automationRepository) BulkUpdateNextRun(ctx context.Context, tasks []*model.RunAutomation) error {
	return r.Executor(ctx).WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "automation_id"}}, // Primary Key
			DoUpdates: clause.AssignmentColumns([]string{
				"status",
				"next_run_time",
				"last_upd",
			}),
		}).
		Create(&tasks).Error
}
