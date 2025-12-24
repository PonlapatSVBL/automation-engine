package repository

import "gorm.io/gorm"

type ConditionOperatorRepository interface {
}

type conditionOperatorRepository struct {
	db *gorm.DB
}

func NewConditionOperatorRepository(db *gorm.DB) ConditionOperatorRepository {
	return &conditionOperatorRepository{db: db}
}
