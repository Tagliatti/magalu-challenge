package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Repository interface {
	CreateNotification(ctx context.Context, createNotification *CreateNotification) (int64, error)
	UpdateNotificationAsSent(ctx context.Context, id int64) (bool, error)
	FindNotificationByID(ctx context.Context, id int64) (*Notification, error)
	FindNotificationStatusByID(ctx context.Context, id int64) (*NotificationStatus, error)
	DeleteNotificationByID(ctx context.Context, id int64) (bool, error)
	ListPendingNotifications(ctx context.Context) ([]Notification, error)
	ListenNotifications(ctx context.Context) (<-chan Notification, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateNotification(ctx context.Context, createNotification *CreateNotification) (int64, error) {
	var id int64

	err := r.db.QueryRow(ctx, `INSERT INTO notifications (type, recipient) VALUES ($1, $2) RETURNING id`,
		createNotification.Type,
		createNotification.Recipient,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresRepository) UpdateNotificationAsSent(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.Exec(ctx, `UPDATE notifications SET sent_at = NOW() WHERE id = $1 and sent_at is null`, id)

	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (r *PostgresRepository) FindNotificationByID(ctx context.Context, id int64) (*Notification, error) {
	var notification Notification
	row := r.db.QueryRow(ctx, `SELECT id, type, recipient, created_at, (sent_at is not null) AS sent, sent_at FROM notifications WHERE id = $1`, id)
	err := row.Scan(
		&notification.Id,
		&notification.Type,
		&notification.Recipient,
		&notification.CreatedAt,
		&notification.Sent,
		&notification.SentAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &notification, nil
}

func (r *PostgresRepository) FindNotificationStatusByID(ctx context.Context, id int64) (*NotificationStatus, error) {
	var notification NotificationStatus
	row := r.db.QueryRow(ctx, `SELECT (sent_at is not null) AS sent, sent_at FROM notifications WHERE id = $1`, id)
	err := row.Scan(
		&notification.Sent,
		&notification.SentAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &notification, nil
}

func (r *PostgresRepository) DeleteNotificationByID(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.Exec(ctx, `DELETE FROM notifications WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (r *PostgresRepository) ListPendingNotifications(ctx context.Context) ([]Notification, error) {
	rows, err := r.db.Query(ctx, `SELECT id, type, recipient, created_at, (sent_at is not null) AS sent, sent_at FROM notifications WHERE sent_at IS NULL`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notifications := make([]Notification, 0)

	for rows.Next() {
		var notification Notification
		err := rows.Scan(
			&notification.Id,
			&notification.Type,
			&notification.Recipient,
			&notification.CreatedAt,
			&notification.Sent,
			&notification.SentAt,
		)

		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *PostgresRepository) ListenNotifications(ctx context.Context) (<-chan Notification, error) {
	channel := make(chan Notification)

	conn, err := r.db.Acquire(ctx)

	if err != nil {
		close(channel)
		return nil, err
	}

	go func() {
		defer close(channel)
		defer conn.Release()

		_, err = conn.Conn().Exec(ctx, `LISTEN pending_notifications`)

		if err != nil {
			log.Printf("Error while listening to notifications: %v", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping notification listening due to context cancellation")
				return
			default:
				pgNotification, err := conn.Conn().WaitForNotification(ctx)
				if err != nil {
					log.Printf("Error while waiting for notification: %v", err)
					return
				}

				var notification Notification

				if err := json.Unmarshal([]byte(pgNotification.Payload), &notification); err != nil {
					log.Printf("Error deserializing notification: %v", err)
					continue
				}

				select {
				case channel <- notification:
				case <-ctx.Done():
					log.Println("Stopping notification channel due to context cancellation")
					return
				}
			}
		}
	}()

	return channel, nil
}
