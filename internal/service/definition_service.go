package service

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/repository"
	"context"
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
	txManager     repository.TransactionManager
	actionRepo    repository.ActionRepository
	conditionRepo repository.ConditionRepository
	operatorRepo  repository.OperatorRepository
	unitRepo      repository.UnitRepository
}

func NewDefinitionService(
	txManager repository.TransactionManager,
	actionRepo repository.ActionRepository,
	conditionRepo repository.ConditionRepository,
	operatorRepo repository.OperatorRepository,
	unitRepo repository.UnitRepository,
) DefinitionService {
	return &definitionService{
		txManager:     txManager,
		actionRepo:    actionRepo,
		conditionRepo: conditionRepo,
		operatorRepo:  operatorRepo,
		unitRepo:      unitRepo,
	}
}

func (s *definitionService) CreateAction(ctx context.Context, action *model.DefAction) error {
	return s.actionRepo.Create(ctx, action)
}

func (s *definitionService) GetActionByID(ctx context.Context, id string) (*model.DefAction, error) {
	return s.actionRepo.GetByID(ctx, id)
}
