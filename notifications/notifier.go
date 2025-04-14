package notifications

import (
	"log"
	"time"
)

const (
	NotificationTypeSms      = "sms"
	NotificationTypeEmail    = "email"
	NotificationTypePush     = "push"
	NotificationTypeWhatsApp = "whatsapp"
)

type Notifier interface {
	Send(notification *Notification) error
}

type NotifierManager struct {
	notifiers map[string]Notifier
}

func (n *NotifierManager) Register(notifier Notifier, notificationType string) {
	n.notifiers[notificationType] = notifier
}

func (n *NotifierManager) Send(notification *Notification) error {
	return n.notifiers[notification.Type].Send(notification)
}

func NewNotifierManager() *NotifierManager {
	return &NotifierManager{
		notifiers: make(map[string]Notifier),
	}
}

type Sms struct{}

func NewSms() *Sms {
	return &Sms{}
}

func (s *Sms) Send(notification *Notification) error {
	log.Printf("Sending SNS notification %d", notification.Id)
	time.Sleep(3 * time.Second) // simulate service delay

	return nil
}

type Email struct{}

func NewEmail() *Email {
	return &Email{}
}

func (e *Email) Send(notification *Notification) error {
	log.Printf("Sending Email notification %d", notification.Id)
	time.Sleep(3 * time.Second) // simulate service delay

	return nil
}

type Push struct{}

func NewPush() *Push {
	return &Push{}
}

func (p *Push) Send(notification *Notification) error {
	log.Printf("Sending Push notification %d", notification.Id)
	time.Sleep(3 * time.Second) // simulate service delay

	return nil
}

type Whatsapp struct{}

func NewWhatsApp() *Whatsapp {
	return &Whatsapp{}
}

func (w *Whatsapp) Send(notification *Notification) error {
	log.Printf("Sending WhatsApp notification %d", notification.Id)
	time.Sleep(3 * time.Second) // simulate service delay

	return nil
}
