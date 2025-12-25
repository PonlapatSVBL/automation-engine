package api

import (
	"automation-engine/internal/service"

	"github.com/gin-gonic/gin"
)

type RunHandler struct {
	runService service.RunService
}

func NewRunHandler(runService service.RunService) *RunHandler {
	return &RunHandler{
		runService: runService,
	}
}

func (h *RunHandler) CreateAutomation(c *gin.Context) {
	// Implementation for creating an automation run
}
