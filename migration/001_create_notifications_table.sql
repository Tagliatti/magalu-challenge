CREATE TYPE notification_type AS ENUM ('email', 'sms', 'push', 'whatsapp');

CREATE TABLE notifications
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type       notification_type NOT NULL,
    recipient  VARCHAR(255)      NOT NULL,
    sent_at    TIMESTAMP         DEFAULT NULL
);
