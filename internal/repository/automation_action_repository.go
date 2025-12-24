package repository

import "gorm.io/gorm"

type AutomationActionRepository interface {
}

type automationActionRepository struct {
	db *gorm.DB
}

func NewAutomationActionRepository(db *gorm.DB) AutomationActionRepository {
	return &automationActionRepository{db: db}
}
