// Code generated manually - equivalent to protoc-gen-go output
// Source: proto/sme/sme.proto

package smepb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── Messages ────────────────────────────────────────────────────────────────

type ListSMEsRequest struct {
	CategoryId string `protobuf:"bytes,1,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"`
	Status     string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Search     string `protobuf:"bytes,3,opt,name=search,proto3" json:"search,omitempty"`
}

type SME struct {
	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId     string   `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	Name        string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Phone       string   `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	Address     string   `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	Description string   `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	CategoryIds []string `protobuf:"bytes,7,rep,name=category_ids,json=categoryIds,proto3" json:"category_ids,omitempty"`
	Products    []string `protobuf:"bytes,8,rep,name=products,proto3" json:"products,omitempty"`
	Capacity    string   `protobuf:"bytes,9,opt,name=capacity,proto3" json:"capacity,omitempty"`
	Latitude    float64  `protobuf:"fixed64,10,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude   float64  `protobuf:"fixed64,11,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Status      string   `protobuf:"bytes,12,opt,name=status,proto3" json:"status,omitempty"`
}

type ListSMEsResponse struct {
	Data  []*SME `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	Total int32  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

// ─── Service Interface ────────────────────────────────────────────────────────

// SMEServiceClient is the client API for SMEService.
type SMEServiceClient interface {
	ListSMEs(ctx context.Context, in *ListSMEsRequest, opts ...grpc.CallOption) (*ListSMEsResponse, error)
}

type sMEServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSMEServiceClient(cc grpc.ClientConnInterface) SMEServiceClient {
	return &sMEServiceClient{cc}
}

func (c *sMEServiceClient) ListSMEs(ctx context.Context, in *ListSMEsRequest, opts ...grpc.CallOption) (*ListSMEsResponse, error) {
	out := new(ListSMEsResponse)
	err := c.cc.Invoke(ctx, "/sme.SMEService/ListSMEs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SMEServiceServer is the server API for SMEService.
type SMEServiceServer interface {
	ListSMEs(context.Context, *ListSMEsRequest) (*ListSMEsResponse, error)
	mustEmbedUnimplementedSMEServiceServer()
}

// UnimplementedSMEServiceServer must be embedded to have forward-compatible implementations.
type UnimplementedSMEServiceServer struct{}

func (UnimplementedSMEServiceServer) ListSMEs(context.Context, *ListSMEsRequest) (*ListSMEsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSMEs not implemented")
}

func (UnimplementedSMEServiceServer) mustEmbedUnimplementedSMEServiceServer() {}

// RegisterSMEServiceServer registers the server to the gRPC server.
func RegisterSMEServiceServer(s grpc.ServiceRegistrar, srv SMEServiceServer) {
	s.RegisterService(&SMEService_ServiceDesc, srv)
}

// SMEService_ServiceDesc is the grpc.ServiceDesc for SMEService.
var SMEService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sme.SMEService",
	HandlerType: (*SMEServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSMEs",
			Handler:    _SMEService_ListSMEs_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func _SMEService_ListSMEs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSMEsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SMEServiceServer).ListSMEs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sme.SMEService/ListSMEs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SMEServiceServer).ListSMEs(ctx, req.(*ListSMEsRequest))
	}
	return interceptor(ctx, in, info, handler)
}
