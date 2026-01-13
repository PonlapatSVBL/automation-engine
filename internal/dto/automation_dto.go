package dto

import "automation-engine/internal/domain/model"

type AutomationSnapshot struct {
	Automation      *model.RunAutomation                 `json:"automation"`
	ConditionGroups []*model.RunAutomationConditionGroup `json:"condition_groups"`
	Conditions      []*model.RunAutomationCondition      `json:"conditions"`
	Actions         []*model.RunAutomationAction         `json:"actions"`
	Targets         []*model.RunAutomationTarget         `json:"targets"`
}
