package repository

import (
	"gorm.io/gorm"
)

type AutomationExecutionRepository interface {
}

type automationExecutionRepository struct {
	BaseRepository
}

func NewAutomationExecutionRepository(db *gorm.DB) AutomationExecutionRepository {
	return &automationExecutionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}
