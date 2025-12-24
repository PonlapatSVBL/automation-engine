package repository

import "gorm.io/gorm"

type UnitRepository interface {
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}
