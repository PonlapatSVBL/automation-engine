package repository

import (
	"context"

	"gorm.io/gorm"
)

type PolicyRepository interface {
}

type policyRepository struct {
	db *gorm.DB
}

type ConditionRow struct {
	ConditionID   string `gorm:"column:condition_id"`
	ConditionCode string `gorm:"column:condition_code"`
	ConditionName string `gorm:"column:condition_name"`
	Status        string `gorm:"column:status"`
}

type ConditionOperatorRow struct {
	ConditionID string  `gorm:"column:condition_id"`
	OperatorID  *string `gorm:"column:operator_id"`
}

type ConditionUnitRow struct {
	ConditionID string  `gorm:"column:condition_id"`
	UnitID      *string `gorm:"column:unit_id"`
}

type ConditionActionRow struct {
	ConditionID string  `gorm:"column:condition_id"`
	ActionID    *string `gorm:"column:action_id"`
}

func NewPolicyRepository(db *gorm.DB) PolicyRepository {
	return &policyRepository{db: db}
}

func (r *policyRepository) GetPolicyRuleConfigs(ctx context.Context) ([]ConditionOperatorRow, []ConditionUnitRow, []ConditionActionRow, error) {
	var conditions []ConditionRow
	var operators []ConditionOperatorRow
	var units []ConditionUnitRow
	var actions []ConditionActionRow

	sqlConditions := `
	SELECT
		_dc.condition_id
		, _dc.condition_code
		, _dc.condition_name
		, _dc.status
	FROM def_conditions _dc
	WHERE _dc.status <> 'inactive'
	`
	if err := r.db.WithContext(ctx).
		Raw(sqlConditions).
		Scan(&conditions).Error; err != nil {
		return nil, nil, nil, err
	}

	sqlOperators := `
	SELECT
		_dc.condition_id
		, _pco.operator_id
	FROM def_conditions _dc
	LEFT JOIN policy_condition_operators _pco
		ON _dc.condition_id = _pco.condition_id
	WHERE _dc.status <> 'inactive'
	`
	if err := r.db.WithContext(ctx).
		Raw(sqlOperators).
		Scan(&operators).Error; err != nil {
		return nil, nil, nil, err
	}

	sqlUnits := `
	SELECT
		_dc.condition_id
		, _pcu.unit_id
	FROM def_conditions _dc
	LEFT JOIN policy_condition_units _pcu
		ON _dc.condition_id = _pcu.condition_id
	WHERE _dc.status <> 'inactive'
	`
	if err := r.db.WithContext(ctx).
		Raw(sqlUnits).
		Scan(&units).Error; err != nil {
		return nil, nil, nil, err
	}

	sqlActions := `
	SELECT
		_dc.condition_id
		, _pca.action_id
	FROM def_conditions _dc
	LEFT JOIN policy_condition_actions _pca
		ON _dc.condition_id = _pca.condition_id
	WHERE _dc.status <> 'inactive'
	`
	if err := r.db.WithContext(ctx).
		Raw(sqlActions).
		Scan(&actions).Error; err != nil {
		return nil, nil, nil, err
	}

	return operators, units, actions, nil
}
