package repository

import "gorm.io/gorm"

type OperatorRepository interface{}

type operatorRepository struct {
	db *gorm.DB
}

func NewOperatorRepository(db *gorm.DB) OperatorRepository {
	return &operatorRepository{db: db}
}
