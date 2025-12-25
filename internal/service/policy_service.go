package service

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/repository"
	"context"

	"gorm.io/gorm"
)

type PolicyService interface {
	GetPolicyRuleConfig(ctx context.Context) (GetPolicyRuleConfigResponse, error)
}

type policyService struct {
	conditionRepo         repository.ConditionRepository
	operatorRepo          repository.OperatorRepository
	unitRepo              repository.UnitRepository
	actionRepo            repository.ActionRepository
	conditionOperatorRepo repository.ConditionOperatorRepository
	conditionUnitRepo     repository.ConditionUnitRepository
	conditionActionRepo   repository.ConditionActionRepository
	// policyRepo            repository.PolicyRepository
}

type GetPolicyRuleConfigResponse struct {
	Conditions         []*model.DefCondition
	Operators          []*model.DefOperator
	Units              []*model.DefUnit
	Actions            []*model.DefAction
	ConditionOperators []*model.PolicyConditionOperator
	ConditionUnits     []*model.PolicyConditionUnit
	ConditionActions   []*model.PolicyConditionAction
}

func NewPolicyService(db *gorm.DB) PolicyService {
	return &policyService{
		conditionRepo:         repository.NewConditionRepository(db),
		operatorRepo:          repository.NewOperatorRepository(db),
		unitRepo:              repository.NewUnitRepository(db),
		actionRepo:            repository.NewActionRepository(db),
		conditionOperatorRepo: repository.NewConditionOperatorRepository(db),
		conditionUnitRepo:     repository.NewConditionUnitRepository(db),
		conditionActionRepo:   repository.NewConditionActionRepository(db),
		// policyRepo:            repository.NewPolicyRepository(db),
	}
}

func (s *policyService) GetPolicyRuleConfig(ctx context.Context) (GetPolicyRuleConfigResponse, error) {
	conditions, err := s.conditionRepo.List(ctx, model.DefCondition{})
	if err != nil {
		return GetPolicyRuleConfigResponse{}, err
	}

	operators, err := s.operatorRepo.List(ctx, model.DefOperator{})
	if err != nil {
		return GetPolicyRuleConfigResponse{}, err
	}

	units, err := s.unitRepo.List(ctx, model.DefUnit{})
	if err != nil {
		return GetPolicyRuleConfigResponse{}, err
	}

	actions, err := s.actionRepo.List(ctx, model.DefAction{})
	if err != nil {
		return GetPolicyRuleConfigResponse{}, err
	}

	conditionOperators, err := s.conditionOperatorRepo.List(ctx, model.PolicyConditionOperator{})
	if err != nil {
		return GetPolicyRuleConfigResponse{}, err
	}

	conditionUnits, err := s.conditionUnitRepo.List(ctx, model.PolicyConditionUnit{})
	if err != nil {
		return GetPolicyRuleConfigResponse{}, err
	}

	conditionActions, err := s.conditionActionRepo.List(ctx, model.PolicyConditionAction{})
	if err != nil {
		return GetPolicyRuleConfigResponse{}, err
	}

	return GetPolicyRuleConfigResponse{
		Conditions:         conditions,
		Operators:          operators,
		Units:              units,
		Actions:            actions,
		ConditionOperators: conditionOperators,
		ConditionUnits:     conditionUnits,
		ConditionActions:   conditionActions,
	}, nil
}
