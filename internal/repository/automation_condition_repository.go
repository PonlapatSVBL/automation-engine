package repository

import "gorm.io/gorm"

type AutomationConditionRepository interface {
}

type automationConditionRepository struct {
	db *gorm.DB
}

func NewAutomationConditionRepository(db *gorm.DB) AutomationConditionRepository {
	return &automationConditionRepository{db: db}
}
