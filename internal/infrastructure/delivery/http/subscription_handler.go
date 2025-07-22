package http

import (
	"net/http"
	"time"

	"github.com/Babushkin05/subscription-organizer/internal/application/port"
	"github.com/Babushkin05/subscription-organizer/internal/shared/dto"
	"github.com/Babushkin05/subscription-organizer/internal/shared/mapper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service port.SubscriptionService
}

func NewSubscriptionHandler(service port.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// POST /subscriptions
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := mapper.ToSubscriptionModel(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use MM-YYYY"})
		return
	}

	err = h.service.CreateSubscription(c.Request.Context(), sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	resp := mapper.ToSubscriptionResponse(*sub)
	c.JSON(http.StatusCreated, resp)
}

// GET /subscriptions/:id
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	sub, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	resp := mapper.ToSubscriptionResponse(*sub)
	c.JSON(http.StatusOK, resp)
}

// PUT /subscriptions/:id
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := mapper.ToSubscriptionModel(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use MM-YYYY"})
		return
	}

	sub.ID = id

	err = h.service.UpdateSubscription(c.Request.Context(), sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
		return
	}

	resp := mapper.ToSubscriptionResponse(*sub)
	c.JSON(http.StatusOK, resp)
}

// DELETE /subscriptions/:id
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	err = h.service.DeleteSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription deleted"})
}

// GET /subscriptions
// Опциональные query-параметры: user_id, service_name
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")

	var userID *uuid.UUID
	if userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = &uid
	}

	// Так как в интерфейсе ListSubscriptions() без фильтров,
	// фильтрация в репозитории или сервисе - нужно расширить, если нужно.
	subs, err := h.service.ListSubscriptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
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

	c.JSON(http.StatusOK, filtered)
}

// GET /subscriptions/cost?user_id=&service_name=&from=&to=
// Подсчет общей стоимости за период с фильтрацией
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'from' and 'to' query parameters required, format MM-YYYY"})
		return
	}

	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date format, use MM-YYYY"})
		return
	}

	to, err := time.Parse("01-2006", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date format, use MM-YYYY"})
		return
	}

	var userID *uuid.UUID
	if userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total cost"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}
