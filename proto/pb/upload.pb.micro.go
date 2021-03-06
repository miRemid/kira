// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/upload.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for UploadService service

func NewUploadServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for UploadService service

type UploadService interface {
	UploadFile(ctx context.Context, in *UploadFileReq, opts ...client.CallOption) (*UploadFileRes, error)
	Ping(ctx context.Context, in *Ping, opts ...client.CallOption) (*Pong, error)
}

type uploadService struct {
	c    client.Client
	name string
}

func NewUploadService(name string, c client.Client) UploadService {
	return &uploadService{
		c:    c,
		name: name,
	}
}

func (c *uploadService) UploadFile(ctx context.Context, in *UploadFileReq, opts ...client.CallOption) (*UploadFileRes, error) {
	req := c.c.NewRequest(c.name, "UploadService.UploadFile", in)
	out := new(UploadFileRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadService) Ping(ctx context.Context, in *Ping, opts ...client.CallOption) (*Pong, error) {
	req := c.c.NewRequest(c.name, "UploadService.Ping", in)
	out := new(Pong)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UploadService service

type UploadServiceHandler interface {
	UploadFile(context.Context, *UploadFileReq, *UploadFileRes) error
	Ping(context.Context, *Ping, *Pong) error
}

func RegisterUploadServiceHandler(s server.Server, hdlr UploadServiceHandler, opts ...server.HandlerOption) error {
	type uploadService interface {
		UploadFile(ctx context.Context, in *UploadFileReq, out *UploadFileRes) error
		Ping(ctx context.Context, in *Ping, out *Pong) error
	}
	type UploadService struct {
		uploadService
	}
	h := &uploadServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UploadService{h}, opts...))
}

type uploadServiceHandler struct {
	UploadServiceHandler
}

func (h *uploadServiceHandler) UploadFile(ctx context.Context, in *UploadFileReq, out *UploadFileRes) error {
	return h.UploadServiceHandler.UploadFile(ctx, in, out)
}

func (h *uploadServiceHandler) Ping(ctx context.Context, in *Ping, out *Pong) error {
	return h.UploadServiceHandler.Ping(ctx, in, out)
}
