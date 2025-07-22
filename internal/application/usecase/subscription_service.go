package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Babushkin05/subscription-organizer/internal/application/port"
	"github.com/Babushkin05/subscription-organizer/internal/domain/model"
	"github.com/google/uuid"
)

type subscriptionService struct {
	repo port.SubscriptionRepository
}

func NewSubscriptionService(repo port.SubscriptionRepository) port.SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, sub *model.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if sub == nil {
		return errors.New("subscription not found")
	}
	sub.IsDeleted = true
	return s.repo.Update(ctx, sub)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	return s.repo.List(ctx)
}

func (s *subscriptionService) CalculateTotalCost(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from time.Time,
	to time.Time,
) (int, error) {
	subs, err := s.repo.GetByFilter(ctx, userID, serviceName, from, to)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, sub := range subs {
		if sub.IsDeleted {
			continue
		}
		months := monthsBetween(sub.StartDate, sub.EndDate, from, to)
		total += sub.Price * months
	}

	return total, nil
}

func monthsBetween(start time.Time, end *time.Time, rangeFrom, rangeTo time.Time) int {
	subStart := maxTime(start, rangeFrom)
	subEnd := rangeTo
	if end != nil && end.Before(rangeTo) {
		subEnd = *end
	}

	if subEnd.Before(subStart) {
		return 0
	}

	// Округляем вниз, чтобы не учитывать неполные месяцы
	years := subEnd.Year() - subStart.Year()
	months := int(subEnd.Month()) - int(subStart.Month()) + years*12 + 1
	return months
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
