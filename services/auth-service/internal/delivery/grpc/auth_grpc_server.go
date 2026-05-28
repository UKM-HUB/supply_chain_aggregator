package grpc

import (
	"context"
	"net"
	"fmt"

	pb "supply-chain-aggregator/services/auth-service/internal/delivery/grpc/pb"
	"supply-chain-aggregator/services/auth-service/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthGRPCServer implements the AuthService gRPC server.
type AuthGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	authUsecase *usecase.AuthUsecase
}

// NewAuthGRPCServer creates a new AuthGRPCServer.
func NewAuthGRPCServer(authUsecase *usecase.AuthUsecase) *AuthGRPCServer {
	return &AuthGRPCServer{authUsecase: authUsecase}
}

// ValidateToken validates a JWT and returns the claims.
func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := s.authUsecase.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{IsValid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		IsValid: true,
		UserId:  claims.UserID,
		Role:    claims.Role,
	}, nil
}

// Start registers the gRPC server and begins listening on grpcPort.
func Start(grpcPort string, authUsecase *usecase.AuthUsecase) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return fmt.Errorf("grpc: failed to listen on port %s: %w", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, NewAuthGRPCServer(authUsecase))

	fmt.Printf("[auth-service] gRPC server listening on :%s\n", grpcPort)
	return grpcServer.Serve(lis)
}
