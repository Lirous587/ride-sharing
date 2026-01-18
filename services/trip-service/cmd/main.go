package main

import (
	"log"
	"net/http"
	handler "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	log.Println("Starting Trip Server")

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	h := handler.HttpHandler{
		Service: svc,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /trip/preview", h.HandlePreviewTrip)

	server := http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
		return
	}
}
