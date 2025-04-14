package main

import (
	"context"
	"github.com/Tagliatti/magalu-challenge/database"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"log"
)

var notifier = notifications.NewNotifierManager()
var repository *notifications.PostgresRepository

func processNotification(ctx context.Context, notification *notifications.Notification) {
	err := notifier.Send(notification)

	if err != nil {
		log.Printf("Error sending notification (%s) %d: %v", notification.Type, notification.Id, err)
		return
	}

	_, err = repository.UpdateNotificationAsSent(ctx, notification.Id)

	if err != nil {
		log.Printf("Error on updating notification %d as sent: %v", notification.Id, err)
		return
	}

	log.Printf("Notification (%s) %d processed successfully", notification.Type, notification.Id)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.Connect(ctx, 10, 10)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	notifier.Register(notifications.NewSms(), notifications.NotificationTypeSms)
	notifier.Register(notifications.NewEmail(), notifications.NotificationTypeEmail)
	notifier.Register(notifications.NewPush(), notifications.NotificationTypePush)
	notifier.Register(notifications.NewWhatsApp(), notifications.NotificationTypeWhatsApp)
	repository = notifications.NewPostgresRepository(db)

	go func() {
		listNotifications, err := repository.ListPendingNotifications(ctx)

		if err != nil {
			log.Fatal(err)
		}

		for _, notification := range listNotifications {
			select {
			case <-ctx.Done():
				log.Println("Stopping pending notifications processing due to context cancellation")
				return
			default:
				processNotification(ctx, &notification)
			}
		}
	}()

	notificationList, err := repository.ListenNotifications(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting for notifications...")
	for notification := range notificationList {
		processNotification(ctx, &notification)
	}
}
