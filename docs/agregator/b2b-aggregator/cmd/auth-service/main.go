package main

import (
	"log"
	"net"
	"b2b-aggregator/internal/auth/delivery/grpc"
	pb "b2b-aggregator/pb/auth"
	grpcServer "google.golang.org/grpc"
)

func main() {
	handler := grpc.NewAuthGrpcServer()
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Gagal listen di port 50053: %v", err)
	}

	s := grpcServer.NewServer()
	pb.RegisterAuthServiceServer(s, handler)

	log.Println("Auth Service (gRPC) berjalan di port 50053...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Gagal serve gRPC: %v", err)
	}
}
