package notifications

import "time"

type Notification struct {
	Id        int64      `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Type      string     `json:"type"`
	Recipient string     `json:"recipient"`
	Sent      bool       `json:"sent"`
	SentAt    *time.Time `json:"sent_at"`
}

type CreateNotification struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
}

type NotificationStatus struct {
	Sent   bool       `json:"sent"`
	SentAt *time.Time `json:"sent_at"`
}
