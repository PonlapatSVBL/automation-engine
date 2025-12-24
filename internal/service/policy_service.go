package service

import (
	"automation-engine/internal/repository"

	"gorm.io/gorm"
)

type PolicyService interface{}

type policyService struct {
	conditionActionRepo   repository.ConditionActionRepository
	conditionOperatorRepo repository.ConditionOperatorRepository
	conditionUnitRepo     repository.ConditionUnitRepository
}

func NewPolicyService(db *gorm.DB) PolicyService {
	return &policyService{
		conditionActionRepo:   repository.NewConditionActionRepository(db),
		conditionOperatorRepo: repository.NewConditionOperatorRepository(db),
		conditionUnitRepo:     repository.NewConditionUnitRepository(db),
	}
}
