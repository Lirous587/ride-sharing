package grpc

import (
	"context"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
	"ride-sharing/shared/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := &types.Coordinate{
		Latitude:  req.StartLocation.GetLatitude(),
		Longitude: req.StartLocation.GetLongitude(),
	}
	destination := &types.Coordinate{
		Latitude:  req.EndLocation.GetLatitude(),
		Longitude: req.EndLocation.GetLongitude(),
	}

	userID := req.GetUserID()

	route, err := h.service.GetRoute(ctx, pickup, destination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	// 1. Estimate the ride fares prices based on the route (ex:distance)
	estimatedFares := h.service.EstimatePackagePriceWithRoute(route)
	// 2. Store the ride fares for the create trip (next soon) to fetch and validate.
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, userID, route)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to generate the ride fares: %v", err)
	}

	protoFares := domain.ToRideFaresProto(fares)

	res := &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: protoFares,
	}

	return res, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()

	// 1.Fetch and validate the fare.
	fare, err := h.service.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("failed to validate fare's legality: %v", err))
	}

	// 2.Call create trip.
	trip, err := h.service.CreateTrip(ctx, fare)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create trip: %v", err))
	}

	// 3.We also need to initialize an empty driver to the trip.
	trip.Driver = &pb.TripDriver{
		Id:             "12304",
		Name:           "fuck damn",
		ProfilePicture: util.GetRandomAvatar(2778),
		CarePlate:      "10086",
	}

	// 4.Add a comment at the end of the function to publish an event on the Async Comms moudle
	// TODO publish an event
	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
