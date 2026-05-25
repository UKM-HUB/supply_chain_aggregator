package grpc

import (
	"context"
	"b2b-aggregator/internal/inventory/usecase"
	pb "b2b-aggregator/pb/inventory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryGrpcServer struct {
	pb.UnimplementedInventoryServiceServer
	usecase *usecase.InventoryUsecase
}

// Injeksi Usecase lewat constructor
func NewInventoryGrpcServer(uc *usecase.InventoryUsecase) *InventoryGrpcServer {
	return &InventoryGrpcServer{usecase: uc}
}

func (s *InventoryGrpcServer) SearchNearby(ctx context.Context, req *pb.SearchNearbyRequest) (*pb.SearchNearbyResponse, error) {
	// Panggil logika bisnis di Usecase
	results, err := s.usecase.SearchNearbySupplies(ctx, req.ProductCode, req.Latitude, req.Longitude, req.RadiusKm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal mencari supply: %v", err)
	}

	// Mapping dari tipe Usecase ke tipe Protobuf
	var pbSupplies []*pb.SupplyItem
	for _, res := range results {
		pbSupplies = append(pbSupplies, &pb.SupplyItem{
			UmkmId:   res.UMKMID,
			Stock:    res.Stock,
			Distance: res.Distance,
		})
	}

	return &pb.SearchNearbyResponse{
		Supplies: pbSupplies,
	}, nil
}
