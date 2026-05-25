package grpc

import (
	"context"
	pb "b2b-aggregator/pb/auth"
	"b2b-aggregator/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcServer struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthGrpcServer() *AuthGrpcServer {
	return &AuthGrpcServer{}
}

func (s *AuthGrpcServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := jwt.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{IsValid: false}, status.Errorf(codes.Unauthenticated, "Token tidak valid")
	}

	return &pb.ValidateTokenResponse{
		IsValid: true,
		UserId:  claims.UserID,
		Role:    claims.Role,
	}, nil
}
