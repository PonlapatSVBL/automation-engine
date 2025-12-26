package api

import (
	"automation-engine/internal/domain/model"
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

type CreateConditionOperatorsRequest struct {
	ConditionID string                           `json:"condition_id" binding:"required"`
	Operators   []*model.PolicyConditionOperator `json:"operators" binding:"required,min=0"`
	CreatedBy   string                           `json:"created_by" binding:"required"`
}

type CreateConditionUnitsRequest struct {
	ConditionID string                       `json:"condition_id" binding:"required"`
	Units       []*model.PolicyConditionUnit `json:"units" binding:"required,min=0"`
	CreatedBy   string                       `json:"created_by" binding:"required"`
}

type CreateConditionActionsRequest struct {
	ConditionID string                         `json:"condition_id" binding:"required"`
	Actions     []*model.PolicyConditionAction `json:"actions" binding:"required,min=0"`
	CreatedBy   string                         `json:"created_by" binding:"required"`
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

func (h *PolicyHandler) CreateConditionOperators(c *gin.Context) {
	var req CreateConditionOperatorsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.policyService.SetConditionOperators(c.Request.Context(), req.ConditionID, req.Operators, req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (h *PolicyHandler) CreateConditionUnits(c *gin.Context) {
	var req CreateConditionUnitsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.policyService.SetConditionUnits(c.Request.Context(), req.ConditionID, req.Units, req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (h *PolicyHandler) CreateConditionActions(c *gin.Context) {
	var req CreateConditionActionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.policyService.SetConditionActions(c.Request.Context(), req.ConditionID, req.Actions, req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}
