package api

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefinitionHandler struct {
	definitionService service.DefinitionService
}

func NewDefinitionHandler(definitionService service.DefinitionService) *DefinitionHandler {
	return &DefinitionHandler{
		definitionService: definitionService,
	}
}

type ActionResponse struct {
	ActionID     string `json:"action_id"`
	ActionCode   string `json:"action_code"`
	ActionName   string `json:"action_name"`
	ActionType   string `json:"action_type"`
	InvokeURL    string `json:"invoke_url"`
	InvokeMethod string `json:"invoke_method"`
	InvokeType   string `json:"invoke_type"`
	Status       string `json:"status"`
}

// GetActionByID godoc
// @Summary      Get action by ID
// @Description  ดึงข้อมูลนิยามของ Action จากตาราง def_actions
// @Tags         definition
// @Accept       json
// @Produce      json
// @Param        id   query      string  true  "Action ID (e.g. ACT001)"
// @Success      200  {object}   api.ActionResponse
// @Failure      400  {object}   map[string]string
// @Failure      404  {object}   map[string]string
// @Router       /definition/actions [get]
// @Security BearerAuth
func (h *DefinitionHandler) GetActionByID(c *gin.Context) {
	// 1. รับค่า ID จาก Path Parameter
	// id := c.Param("id")
	id := c.Query("id")
	fmt.Println(id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action id is required"})
		return
	}

	// 2. เรียกใช้ Service
	// หมายเหตุ: c.Request.Context() ช่วยให้เรายกเลิกการทำงานได้ถ้า Client ตัดการเชื่อมต่อ
	action, err := h.definitionService.GetActionByID(c.Request.Context(), id)
	if err != nil {
		// คุณสามารถแยก error ได้ว่าถ้าไม่เจอให้ส่ง 404
		c.JSON(http.StatusNotFound, gin.H{"error": "action not found"})
		return
	}

	// 3. ส่งข้อมูลกลับ
	c.JSON(http.StatusOK, action)
}

// CreateAction godoc
// @Summary      Create action
// @Description  สร้างนิยาม Action ใหม่
// @Tags         definition
// @Accept       json
// @Produce      json
// @Param        body  body      api.CreateActionRequest  true  "Create Action Payload"
// @Success      201   {object}  api.ActionResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /definition/actions [post]
// @Security BearerAuth
type CreateActionRequest struct {
	ActionID     string `json:"action_id" binding:"required"`
	ActionCode   string `json:"action_code" binding:"required"`
	ActionName   string `json:"action_name" binding:"required"`
	ActionType   string `json:"action_type" binding:"required"`
	InvokeURL    string `json:"invoke_url" binding:"required,url"`
	InvokeMethod string `json:"invoke_method" binding:"required,oneof=GET POST PUT DELETE"`
	InvokeType   string `json:"invoke_type" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

func (h *DefinitionHandler) CreateAction(c *gin.Context) {
	var req CreateActionRequest

	// 1. Bind + Validate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 2. Map request → domain model
	action := &model.DefAction{
		ActionID:     req.ActionID,
		ActionCode:   req.ActionCode,
		ActionName:   req.ActionName,
		ActionType:   req.ActionType,
		InvokeURL:    req.InvokeURL,
		InvokeMethod: req.InvokeMethod,
		InvokeType:   req.InvokeType,
		Status:       req.Status,
	}

	// 3. Call service
	if err := h.definitionService.CreateAction(c.Request.Context(), action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create action",
		})
		return
	}

	// 4. Response
	resp := ActionResponse{
		ActionID:     action.ActionID,
		ActionCode:   action.ActionCode,
		ActionName:   action.ActionName,
		ActionType:   action.ActionType,
		InvokeURL:    action.InvokeURL,
		InvokeMethod: action.InvokeMethod,
		InvokeType:   action.InvokeType,
		Status:       action.Status,
	}

	c.JSON(http.StatusCreated, resp)
}
