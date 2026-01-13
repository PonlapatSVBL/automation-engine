package dto

import (
	"errors"
	"time"
)

type MessageServiceBus struct {
	LogID        string    `json:"log_id" validate:"required"`
	AutomationID string    `json:"automation_id" validate:"required"`
	TriggeredAt  time.Time `json:"triggered_at" validate:"required"`
}

func (m *MessageServiceBus) Validate() error {
	if m.LogID == "" {
		return errors.New("log_id is required")
	}
	if m.AutomationID == "" {
		return errors.New("automation_id is required")
	}
	if m.TriggeredAt.IsZero() {
		return errors.New("triggered_at is required")
	}
	return nil
}
