package main

import (
	"github.com/Tagliatti/magalu-challenge/database"
	"github.com/Tagliatti/magalu-challenge/health"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"github.com/Tagliatti/magalu-challenge/notifications/handler"
	"log"
	"net/http"
)

func main() {
	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	notificationStorage := notifications.NewPostgresRepository(db)

	healthy := health.NewHealthyHandler()
	createNotification := handler.NewCreateHandler(notificationStorage)

	server := http.NewServeMux()
	server.HandleFunc("POST /notifications", createNotification.Handler)
	server.HandleFunc("/", healthy.Handler)

	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", server))
}
