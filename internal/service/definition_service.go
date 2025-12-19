package service

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/repository"
	"context"

	"gorm.io/gorm"
)

type DefinitionService interface {
	// Group CRUD
	// CreateGroup(ctx context.Context, group *model.DefGroup) error
	// GetGroupByID(ctx context.Context, id string) (*model.DefGroup, error)
	// UpdateGroup(ctx context.Context, group *model.DefGroup) error
	// DeleteGroup(ctx context.Context, id string) error

	// Condition CRUD
	// CreateCondition(ctx context.Context, condition *model.DefCondition) error
	// GetConditionByID(ctx context.Context, id string) (*model.DefCondition, error)
	// UpdateCondition(ctx context.Context, condition *model.DefCondition) error
	// DeleteCondition(ctx context.Context, id string) error

	// Action CRUD
	CreateAction(ctx context.Context, action *model.DefAction) error
	GetActionByID(ctx context.Context, id string) (*model.DefAction, error)
	// UpdateAction(ctx context.Context, action *model.DefAction) error
	// DeleteAction(ctx context.Context, id string) error
}

type definitionService struct {
	actionRepo repository.ActionRepository
}

func NewDefinitionService(db *gorm.DB) DefinitionService {
	return &definitionService{
		actionRepo: repository.NewActionRepository(db),
	}
}

func (s *definitionService) CreateAction(ctx context.Context, action *model.DefAction) error {
	return s.actionRepo.Create(ctx, action)
}

func (s *definitionService) GetActionByID(ctx context.Context, id string) (*model.DefAction, error) {
	return s.actionRepo.GetByID(ctx, id)
}
