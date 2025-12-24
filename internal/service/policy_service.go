package service

import (
	"automation-engine/internal/repository"

	"gorm.io/gorm"
)

type PolicyService interface{}

type policyService struct {
	conditionActionRepo repository.ConditionActionRepository
}

func NewPolicyService(db *gorm.DB) PolicyService {
	return &policyService{
		conditionActionRepo: repository.NewConditionActionRepository(db),
	}
}
