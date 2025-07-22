package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Babushkin05/subscription-organizer/internal/application/port"
	"github.com/Babushkin05/subscription-organizer/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) port.SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *model.Subscription) error {
	query := `
		INSERT INTO subscriptions 
		(id, service_name, price, user_id, start_date, end_date, is_deleted)
		VALUES (:id, :service_name, :price, :user_id, :start_date, :end_date, :is_deleted)
	`

	_, err := r.db.NamedExecContext(ctx, query, sub)
	return err
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, is_deleted
		FROM subscriptions
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &sub, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &sub, err
}

func (r *subscriptionRepo) Update(ctx context.Context, sub *model.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = :service_name,
			price = :price,
			user_id = :user_id,
			start_date = :start_date,
			end_date = :end_date,
			is_deleted = :is_deleted
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, sub)
	return err
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE subscriptions
		SET is_deleted = true
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *subscriptionRepo) List(ctx context.Context) ([]*model.Subscription, error) {
	var subs []*model.Subscription

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, is_deleted
		FROM subscriptions
		WHERE is_deleted = false
	`

	err := r.db.SelectContext(ctx, &subs, query)
	return subs, err
}

func (r *subscriptionRepo) GetByFilter(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from time.Time,
	to time.Time,
) ([]*model.Subscription, error) {
	var subs []*model.Subscription

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, is_deleted
		FROM subscriptions
		WHERE is_deleted = false
		  AND start_date <= $1
		  AND (end_date IS NULL OR end_date >= $2)
	`

	args := []interface{}{to, from}

	if userID != nil {
		query += " AND user_id = $3"
		args = append(args, *userID)
		if serviceName != nil {
			query += " AND service_name = $4"
			args = append(args, *serviceName)
		}
	} else if serviceName != nil {
		query += " AND service_name = $3"
		args = append(args, *serviceName)
	}

	err := r.db.SelectContext(ctx, &subs, query, args...)
	return subs, err
}
