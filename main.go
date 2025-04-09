package main

import (
	"github.com/Tagliatti/magalu-challenge/health"
	"log"
	"net/http"
)

func main() {
	healthy := health.NewHealthyHandler()

	server := http.NewServeMux()
	server.HandleFunc("/", healthy.Handler)

	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", server))
}
