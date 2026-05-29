package grpc

import (
	"context"
	"fmt"
	"net"

	pb "supply-chain-aggregator/services/transaction-service/internal/delivery/grpc/pb"
	"supply-chain-aggregator/services/transaction-service/internal/entity"
	"supply-chain-aggregator/services/transaction-service/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransactionGRPCServer implements the TransactionService gRPC server.
type TransactionGRPCServer struct {
	pb.UnimplementedTransactionServiceServer
	txUsecase *usecase.TransactionUsecase
}

func NewTransactionGRPCServer(txUsecase *usecase.TransactionUsecase) *TransactionGRPCServer {
	return &TransactionGRPCServer{txUsecase: txUsecase}
}

func (s *TransactionGRPCServer) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error) {
	tx, err := s.txUsecase.Create(ctx, usecase.CreateTransactionInput{
		UserID:        req.UserId,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create transaction: %v", err)
	}
	return &pb.TransactionResponse{Data: txToProto(tx)}, nil
}

func (s *TransactionGRPCServer) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.TransactionResponse, error) {
	tx, err := s.txUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "transaction not found: %v", err)
	}
	return &pb.TransactionResponse{Data: txToProto(tx)}, nil
}

func (s *TransactionGRPCServer) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	txs, err := s.txUsecase.List(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list transactions: %v", err)
	}
	pbTxs := make([]*pb.Transaction, 0, len(txs))
	for _, tx := range txs {
		pbTxs = append(pbTxs, txToProto(tx))
	}
	return &pb.ListTransactionsResponse{Data: pbTxs, Total: int32(len(pbTxs))}, nil
}

func (s *TransactionGRPCServer) UpdateTransactionStatus(ctx context.Context, req *pb.UpdateTransactionStatusRequest) (*pb.UpdateTransactionStatusResponse, error) {
	if err := s.txUsecase.UpdateStatus(ctx, req.Id, req.Status); err != nil {
		return &pb.UpdateTransactionStatusResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.UpdateTransactionStatusResponse{Success: true, Message: "status updated successfully"}, nil
}

func txToProto(tx entity.Transaction) *pb.Transaction {
	return &pb.Transaction{
		Id:            tx.ID,
		InvoiceNumber: tx.InvoiceNumber,
		UserId:        tx.UserID,
		Amount:        tx.Amount,
		Status:        tx.Status,
		PaymentMethod: tx.PaymentMethod,
		CreatedAt:     tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     tx.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Start registers the gRPC server and begins listening on grpcPort.
func Start(grpcPort string, txUsecase *usecase.TransactionUsecase) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return fmt.Errorf("grpc: failed to listen on port %s: %w", grpcPort, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTransactionServiceServer(grpcServer, NewTransactionGRPCServer(txUsecase))
	fmt.Printf("[transaction-service] gRPC server listening on :%s\n", grpcPort)
	return grpcServer.Serve(lis)
}
