package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var (
	grpcAddr    = env.GetString("HTTP_ADDR", ":9093")
	rabbitMQURI = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMQURI)
	if err != nil {
		log.Fatalln(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	publisher := events.NewTripEventPublisher(rabbitmq)
	driverConsumer := events.NewDriverConsumer(rabbitmq, svc)
	go driverConsumer.Listen()

	// Starting the grpc server
	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, svc, publisher)

	log.Printf("Starting gRPC server Trip service on port: %v", lis.Addr())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
