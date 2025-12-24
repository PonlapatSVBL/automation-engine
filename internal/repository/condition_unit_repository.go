package repository

import "gorm.io/gorm"

type ConditionUnitRepository interface {
}

type conditionUnitRepository struct {
	db *gorm.DB
}

func NewConditionUnitRepository(db *gorm.DB) ConditionUnitRepository {
	return &conditionUnitRepository{db: db}
}
