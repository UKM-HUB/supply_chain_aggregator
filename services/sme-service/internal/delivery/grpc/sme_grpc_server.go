package grpc

import (
	"context"
	"fmt"
	"net"

	pb "supply-chain-aggregator/services/sme-service/internal/delivery/grpc/pb"
	"supply-chain-aggregator/services/sme-service/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SMEGRPCServer implements the SMEService gRPC server.
type SMEGRPCServer struct {
	pb.UnimplementedSMEServiceServer
	smeUsecase *usecase.SMEUsecase
}

// NewSMEGRPCServer creates a new SMEGRPCServer.
func NewSMEGRPCServer(smeUsecase *usecase.SMEUsecase) *SMEGRPCServer {
	return &SMEGRPCServer{smeUsecase: smeUsecase}
}

// ListSMEs returns a filtered list of SMEs via gRPC.
func (s *SMEGRPCServer) ListSMEs(ctx context.Context, req *pb.ListSMEsRequest) (*pb.ListSMEsResponse, error) {
	smes, err := s.smeUsecase.List(ctx, usecase.ListSMEInput{
		CategoryID: req.CategoryId,
		Status:     req.Status,
		Search:     req.Search,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list SMEs: %v", err)
	}

	pbSMEs := make([]*pb.SME, 0, len(smes))
	for _, sme := range smes {
		pbSMEs = append(pbSMEs, &pb.SME{
			Id:          sme.ID,
			OwnerId:     sme.OwnerID,
			Name:        sme.Name,
			Phone:       sme.Phone,
			Address:     sme.Address,
			Description: sme.Description,
			CategoryIds: sme.CategoryIDs,
			Products:    sme.Products,
			Capacity:    sme.Capacity,
			Latitude:    sme.Latitude,
			Longitude:   sme.Longitude,
			Status:      sme.Status,
		})
	}

	return &pb.ListSMEsResponse{
		Data:  pbSMEs,
		Total: int32(len(pbSMEs)),
	}, nil
}

// Start registers the gRPC server and listens on grpcPort.
func Start(grpcPort string, smeUsecase *usecase.SMEUsecase) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return fmt.Errorf("grpc: failed to listen on port %s: %w", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSMEServiceServer(grpcServer, NewSMEGRPCServer(smeUsecase))

	fmt.Printf("[sme-service] gRPC server listening on :%s\n", grpcPort)
	return grpcServer.Serve(lis)
}
