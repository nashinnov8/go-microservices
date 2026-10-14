package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nashinnov8/go-ecommerce/pkg/apperror"
	"github.com/nashinnov8/go-ecommerce/pkg/response"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/dto"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/service"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// In production with JWT, user_id comes from c.GetString("user_id")
	// Fallback to query param or default for dev testing
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		userIDStr = c.Query("user_id")
	}
	if userIDStr == "" {
		userIDStr = "00000000-0000-0000-0000-000000000001" // Default test user
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperror.BadRequest("INVALID_USER_ID", "Invalid user ID format", err))
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res, err := h.svc.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, res)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	raw := c.Param("id")
	id, err := uuid.Parse(raw)
	if err != nil {
		response.Error(c, apperror.BadRequest("INVALID_ORDER_ID", "Invalid order ID format", err))
		return
	}

	res, err := h.svc.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, res)
}
