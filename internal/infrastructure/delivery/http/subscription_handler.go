package http

import (
	"net/http"
	"time"

	"github.com/Babushkin05/subscription-organizer/internal/application/port"
	"github.com/Babushkin05/subscription-organizer/internal/shared/dto"
	"github.com/Babushkin05/subscription-organizer/internal/shared/mapper"
	"github.com/Babushkin05/subscription-organizer/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service port.SubscriptionService
}

func NewSubscriptionHandler(service port.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Creates a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequest true "Subscription to create"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	logger.Log.Info("CreateSubscription: received request")
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sub, err := mapper.ToSubscriptionModel(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid date format, use MM-YYYY"})
		return
	}

	err = h.service.CreateSubscription(c.Request.Context(), sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create subscription"})
		return
	}

	resp := mapper.ToSubscriptionResponse(*sub)
	logger.Log.Infof("CreateSubscription: subscription created with ID %s", resp.ID)
	c.JSON(http.StatusCreated, resp)
}

// GetSubscription godoc
// @Summary Get a subscription by ID
// @Description Returns a single subscription
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	logger.Log.Infof("GetSubscription: getting subscription %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid subscription id"})
		return
	}

	sub, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "subscription not found"})
		return
	}

	resp := mapper.ToSubscriptionResponse(*sub)
	logger.Log.Infof("GetSubscription: found subscription %s", id)
	c.JSON(http.StatusOK, resp)
}

// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Updates a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body dto.CreateSubscriptionRequest true "Updated subscription data"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	logger.Log.Infof("UpdateSubscription: updating subscription %s", idStr)
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid subscription id"})
		return
	}

	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sub, err := mapper.ToSubscriptionModel(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid date format, use MM-YYYY"})
		return
	}

	sub.ID = id

	err = h.service.UpdateSubscription(c.Request.Context(), sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update subscription"})
		return
	}

	resp := mapper.ToSubscriptionResponse(*sub)
	logger.Log.Infof("UpdateSubscription: updated subscription %s", sub.ID)
	c.JSON(http.StatusOK, resp)
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Soft-deletes a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	logger.Log.Infof("DeleteSubscription: deleting subscription %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid subscription id"})
		return
	}

	err = h.service.DeleteSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to delete subscription"})
		return
	}

	logger.Log.Infof("DeleteSubscription: deleted subscription %s", id)
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "subscription deleted"})
}

// ListSubscriptions godoc
// @Summary List all subscriptions
// @Description Returns a list of subscriptions (optionally filtered by user_id and service_name)
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User UUID"
// @Param service_name query string false "Service Name"
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	logger.Log.Infof("ListSubscriptions: query user_id=%s, service_name=%s", userIDStr, serviceName)

	var userID *uuid.UUID
	if userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
			return
		}
		userID = &uid
	}

	subs, err := h.service.ListSubscriptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to list subscriptions"})
		return
	}

	filtered := make([]dto.SubscriptionResponse, 0, len(subs))
	for _, s := range subs {
		if userID != nil && s.UserID != *userID {
			continue
		}
		if serviceName != "" && s.ServiceName != serviceName {
			continue
		}
		filtered = append(filtered, mapper.ToSubscriptionResponse(*s))
	}

	logger.Log.Infof("ListSubscriptions: returned %d subscriptions", len(filtered))
	c.JSON(http.StatusOK, filtered)
}

// CalculateTotalCost godoc
// @Summary Calculate total subscription cost
// @Description Calculates the total cost of subscriptions over a time period with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User UUID"
// @Param service_name query string false "Service Name"
// @Param from query string true "Start period in MM-YYYY"
// @Param to query string true "End period in MM-YYYY"
// @Success 200 {object} map[string]int
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/cost [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	logger.Log.Infof("CalculateTotalCost: user_id=%s, service_name=%s, from=%s, to=%s", userIDStr, serviceName, fromStr, toStr)

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "'from' and 'to' query parameters required, format MM-YYYY"})
		return
	}

	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid 'from' date format, use MM-YYYY"})
		return
	}

	to, err := time.Parse("01-2006", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid 'to' date format, use MM-YYYY"})
		return
	}

	var userID *uuid.UUID
	if userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
			return
		}
		userID = &uid
	}

	var svcNamePtr *string
	if serviceName != "" {
		svcNamePtr = &serviceName
	}

	total, err := h.service.CalculateTotalCost(c.Request.Context(), userID, svcNamePtr, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to calculate total cost"})
		return
	}

	logger.Log.Infof("CalculateTotalCost: total cost = %d", total)
	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}
