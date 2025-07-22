package port

import (
	"context"
	"time"

	"github.com/Babushkin05/subscription-organizer/internal/domain/model"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *model.Subscription) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListSubscriptions(ctx context.Context) ([]*model.Subscription, error)

	CalculateTotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error)
}
