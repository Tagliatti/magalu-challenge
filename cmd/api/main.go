package main

import (
	"context"
	"github.com/Tagliatti/magalu-challenge/database"
	"github.com/Tagliatti/magalu-challenge/health"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"github.com/Tagliatti/magalu-challenge/notifications/handler"
	"log"
	"net"
	"net/http"
)

func main() {
	ctx := context.Background()
	db, err := database.Connect(ctx, 10, 100)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	notificationStorage := notifications.NewPostgresRepository(db)

	healthy := health.NewHealthyHandler()
	createNotification := handler.NewCreateHandler(notificationStorage)
	statusNotification := handler.NewStatusHandler(notificationStorage)
	deleteNotification := handler.NewDeleteHandler(notificationStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /notifications", createNotification.Handler)
	mux.HandleFunc("GET /notifications/{id}/status", statusNotification.Handler)
	mux.HandleFunc("DELETE /notifications/{id}", deleteNotification.Handler)
	mux.HandleFunc("/", healthy.Handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(server.ListenAndServe())
}
