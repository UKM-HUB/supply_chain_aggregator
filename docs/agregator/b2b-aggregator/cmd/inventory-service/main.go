package main

import (
	"log"
	"net"
	"b2b-aggregator/internal/inventory/delivery/grpc"
	"b2b-aggregator/internal/inventory/repository"
	"b2b-aggregator/internal/inventory/usecase"
	pb "b2b-aggregator/pb/inventory"
	"github.com/redis/go-redis/v9"
	grpcServer "google.golang.org/grpc"
)

func main() {
	// 1. Setup Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 2. Setup Clean Architecture Layers
	geoRepo := repository.NewGeoRepository(rdb)
	inventoryUC := usecase.NewInventoryUsecase(geoRepo)
	handler := grpc.NewInventoryGrpcServer(inventoryUC)

	// 3. Jalankan Server gRPC
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Gagal listen di port 50052: %v", err)
	}

	s := grpcServer.NewServer()
	pb.RegisterInventoryServiceServer(s, handler)

	log.Println("Inventory Service (gRPC) berjalan di port 50052...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Gagal serve gRPC: %v", err)
	}
}
