package repository

import "gorm.io/gorm"

type AutomationRepository interface {
}

type automationRepository struct {
	db *gorm.DB
}

func NewAutomationRepository(db *gorm.DB) AutomationRepository {
	return &automationRepository{db: db}
}
