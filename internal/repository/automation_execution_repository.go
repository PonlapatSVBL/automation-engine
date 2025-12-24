package repository

import "gorm.io/gorm"

type AutomationExecutionRepository interface {
}

type automationExecutionRepository struct {
	db *gorm.DB
}

func NewAutomationExecutionRepository(db *gorm.DB) AutomationExecutionRepository {
	return &automationExecutionRepository{db: db}
}
