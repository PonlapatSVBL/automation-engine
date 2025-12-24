package repository

import "gorm.io/gorm"

type ConditionRepository interface{}

type conditionRepository struct {
	db *gorm.DB
}

func NewConditionRepository(db *gorm.DB) ConditionRepository {
	return &conditionRepository{db: db}
}
