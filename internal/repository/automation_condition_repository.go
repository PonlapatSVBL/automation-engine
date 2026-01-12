package repository

import (
	"gorm.io/gorm"
)

type AutomationConditionRepository interface {
}

type automationConditionRepository struct {
	BaseRepository
}

func NewAutomationConditionRepository(db *gorm.DB) AutomationConditionRepository {
	return &automationConditionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}
