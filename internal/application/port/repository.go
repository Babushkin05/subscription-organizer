package port

import (
	"context"
	"time"

	"github.com/Babushkin05/subscription-organizer/internal/domain/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*model.Subscription, error)

	GetByFilter(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) ([]*model.Subscription, error)
}
