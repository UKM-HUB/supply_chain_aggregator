package grpc

import (
	"context"
	"b2b-aggregator/internal/order/usecase"
	pb "b2b-aggregator/pb/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcServer struct {
	pb.UnimplementedOrderServiceServer
	usecase *usecase.OrderUsecase
}

func NewOrderGrpcServer(uc *usecase.OrderUsecase) *OrderGrpcServer {
	return &OrderGrpcServer{usecase: uc}
}

func (s *OrderGrpcServer) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	err := s.usecase.ProcessCheckout(ctx, req.FactoryId, req.ProductCode, req.Quantity)
	if err != nil {
		// Ubah error internal menjadi format gRPC Status
		return nil, status.Errorf(codes.Internal, "gagal memproses checkout: %v", err)
	}

	return &pb.CheckoutResponse{
		Status:  "SUCCESS",
		Message: "Pesanan berhasil dipecah dan disimpan",
	}, nil
}
