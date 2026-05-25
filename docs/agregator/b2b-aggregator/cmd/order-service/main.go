package main

import (
	"log"
	"net"
	"b2b-aggregator/internal/order/delivery/grpc"
	"b2b-aggregator/internal/order/repository"
	"b2b-aggregator/internal/order/usecase"
	"b2b-aggregator/pkg/config"
	"b2b-aggregator/pkg/rabbitmq"
	pb "b2b-aggregator/pb/order"
	grpcServer "google.golang.org/grpc"
)

func main() {
	// 1. Setup DB
	db, err := config.InitDatabase("localhost", "root", "secretpassword", "b2b_db", "5432")
	if err != nil {
		log.Fatalf("Gagal connect DB: %v", err)
	}

	// 2. Setup RabbitMQ
	pub, err := rabbitmq.NewPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("Warning: RabbitMQ tidak terhubung: %v", err)
	}

	// 3. Setup Layer Clean Architecture
	repo := repository.NewOrderRepository(db)
	uc := usecase.NewOrderUsecase(repo, pub)
	handler := grpc.NewOrderGrpcServer(uc)

	// 4. Jalankan gRPC Server di port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal listen di port 50051: %v", err)
	}

	s := grpcServer.NewServer()
	pb.RegisterOrderServiceServer(s, handler)

	log.Println("Order Service (gRPC) berjalan di port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Gagal serve gRPC: %v", err)
	}
}
