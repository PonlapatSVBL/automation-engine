package api

import (
	"automation-engine/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PolicyHandler struct {
	policyService service.PolicyService
}

func NewPolicyHandler(policyService service.PolicyService) *PolicyHandler {
	return &PolicyHandler{
		policyService: policyService,
	}
}

type CreatePolicyConditionActionRequest struct {
	ConditionID string   `json:"condition_id" binding:"required"`
	ActionIDs   []string `json:"action_ids" binding:"required,min=1"`
	CreatedBy   string   `json:"created_by" binding:"required"`
}

func (h *PolicyHandler) GetPolicyRuleConfig(c *gin.Context) {
	response, err := h.policyService.GetPolicyRuleConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *PolicyHandler) CreateConditionActions(c *gin.Context) {
	var req CreatePolicyConditionActionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
}
