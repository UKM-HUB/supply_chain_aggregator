// Code generated manually - equivalent to protoc-gen-go output
// Source: proto/transaction/transaction.proto

package transactionpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── Messages ────────────────────────────────────────────────────────────────

type CreateTransactionRequest struct {
	UserId        string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount        float64 `protobuf:"fixed64,2,opt,name=amount,proto3" json:"amount,omitempty"`
	PaymentMethod string  `protobuf:"bytes,3,opt,name=payment_method,json=paymentMethod,proto3" json:"payment_method,omitempty"`
}

type GetTransactionRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

type ListTransactionsRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type UpdateTransactionStatusRequest struct {
	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

type Transaction struct {
	Id            string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	InvoiceNumber string  `protobuf:"bytes,2,opt,name=invoice_number,json=invoiceNumber,proto3" json:"invoice_number,omitempty"`
	UserId        string  `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount        float64 `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Status        string  `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	PaymentMethod string  `protobuf:"bytes,6,opt,name=payment_method,json=paymentMethod,proto3" json:"payment_method,omitempty"`
	CreatedAt     string  `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     string  `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

type TransactionResponse struct {
	Data *Transaction `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

type ListTransactionsResponse struct {
	Data  []*Transaction `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	Total int32          `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

type UpdateTransactionStatusResponse struct {
	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

// ─── Service Interface ────────────────────────────────────────────────────────

type TransactionServiceClient interface {
	CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error)
	GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error)
	ListTransactions(ctx context.Context, in *ListTransactionsRequest, opts ...grpc.CallOption) (*ListTransactionsResponse, error)
	UpdateTransactionStatus(ctx context.Context, in *UpdateTransactionStatusRequest, opts ...grpc.CallOption) (*UpdateTransactionStatusResponse, error)
}

type transactionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionServiceClient(cc grpc.ClientConnInterface) TransactionServiceClient {
	return &transactionServiceClient{cc}
}

func (c *transactionServiceClient) CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/transaction.TransactionService/CreateTransaction", in, out, opts...)
	return out, err
}

func (c *transactionServiceClient) GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/transaction.TransactionService/GetTransaction", in, out, opts...)
	return out, err
}

func (c *transactionServiceClient) ListTransactions(ctx context.Context, in *ListTransactionsRequest, opts ...grpc.CallOption) (*ListTransactionsResponse, error) {
	out := new(ListTransactionsResponse)
	err := c.cc.Invoke(ctx, "/transaction.TransactionService/ListTransactions", in, out, opts...)
	return out, err
}

func (c *transactionServiceClient) UpdateTransactionStatus(ctx context.Context, in *UpdateTransactionStatusRequest, opts ...grpc.CallOption) (*UpdateTransactionStatusResponse, error) {
	out := new(UpdateTransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/transaction.TransactionService/UpdateTransactionStatus", in, out, opts...)
	return out, err
}

// TransactionServiceServer is the server API for TransactionService.
type TransactionServiceServer interface {
	CreateTransaction(context.Context, *CreateTransactionRequest) (*TransactionResponse, error)
	GetTransaction(context.Context, *GetTransactionRequest) (*TransactionResponse, error)
	ListTransactions(context.Context, *ListTransactionsRequest) (*ListTransactionsResponse, error)
	UpdateTransactionStatus(context.Context, *UpdateTransactionStatusRequest) (*UpdateTransactionStatusResponse, error)
	mustEmbedUnimplementedTransactionServiceServer()
}

type UnimplementedTransactionServiceServer struct{}

func (UnimplementedTransactionServiceServer) CreateTransaction(context.Context, *CreateTransactionRequest) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransaction not implemented")
}
func (UnimplementedTransactionServiceServer) GetTransaction(context.Context, *GetTransactionRequest) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransaction not implemented")
}
func (UnimplementedTransactionServiceServer) ListTransactions(context.Context, *ListTransactionsRequest) (*ListTransactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTransactions not implemented")
}
func (UnimplementedTransactionServiceServer) UpdateTransactionStatus(context.Context, *UpdateTransactionStatusRequest) (*UpdateTransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTransactionStatus not implemented")
}
func (UnimplementedTransactionServiceServer) mustEmbedUnimplementedTransactionServiceServer() {}

// RegisterTransactionServiceServer registers the server.
func RegisterTransactionServiceServer(s grpc.ServiceRegistrar, srv TransactionServiceServer) {
	s.RegisterService(&TransactionService_ServiceDesc, srv)
}

var TransactionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transaction.TransactionService",
	HandlerType: (*TransactionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateTransaction", Handler: _TransactionService_CreateTransaction_Handler},
		{MethodName: "GetTransaction", Handler: _TransactionService_GetTransaction_Handler},
		{MethodName: "ListTransactions", Handler: _TransactionService_ListTransactions_Handler},
		{MethodName: "UpdateTransactionStatus", Handler: _TransactionService_UpdateTransactionStatus_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _TransactionService_CreateTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).CreateTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/CreateTransaction"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).CreateTransaction(ctx, req.(*CreateTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionService_GetTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).GetTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/GetTransaction"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).GetTransaction(ctx, req.(*GetTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionService_ListTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTransactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).ListTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/ListTransactions"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).ListTransactions(ctx, req.(*ListTransactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionService_UpdateTransactionStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTransactionStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).UpdateTransactionStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/UpdateTransactionStatus"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).UpdateTransactionStatus(ctx, req.(*UpdateTransactionStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}
