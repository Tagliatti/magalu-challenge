package notifications

import (
	"database/sql"
	"errors"
)

type Repository interface {
	CreateNotification(createNotification *CreateNotification) (int64, error)
	UpdateNotificationAsSent(id int64) (bool, error)
	FindNotificationByID(id int64) (*Notification, error)
	FindNotificationStatusByID(id int64) (*NotificationStatus, error)
	DeleteNotificationByID(id int64) (bool, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateNotification(createNotification *CreateNotification) (int64, error) {
	var id int64

	err := r.db.QueryRow(`INSERT INTO notifications (type, recipient) VALUES ($1, $2) RETURNING id`,
		createNotification.Type,
		createNotification.Recipient,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresRepository) UpdateNotificationAsSent(id int64) (bool, error) {
	result, err := r.db.Exec(`UPDATE notifications SET sent_at = NOW() WHERE id = $1 and sent_at is null`, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *PostgresRepository) FindNotificationByID(id int64) (*Notification, error) {
	var notification Notification
	row := r.db.QueryRow(`SELECT id, type, recipient, created_at, (sent_at is not null) AS sent, sent_at FROM notifications WHERE id = $1`, id)
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

func (r *PostgresRepository) FindNotificationStatusByID(id int64) (*NotificationStatus, error) {
	var notification NotificationStatus
	row := r.db.QueryRow(`SELECT (sent_at is not null) AS sent, sent_at FROM notifications WHERE id = $1`, id)
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

func (r *PostgresRepository) DeleteNotificationByID(id int64) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM notifications WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
