package repository

import (
	"gorm.io/gorm"
)

type AutomationActionRepository interface {
}

type automationActionRepository struct {
	BaseRepository
}

func NewAutomationActionRepository(db *gorm.DB) AutomationActionRepository {
	return &automationActionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}
