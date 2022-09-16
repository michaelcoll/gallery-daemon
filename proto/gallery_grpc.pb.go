// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: proto/gallery.proto

package gallery

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GalleryClient is the client API for Gallery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GalleryClient interface {
	// get all photos
	GetPhotos(ctx context.Context, in *PhotoFilter, opts ...grpc.CallOption) (Gallery_GetPhotosClient, error)
}

type galleryClient struct {
	cc grpc.ClientConnInterface
}

func NewGalleryClient(cc grpc.ClientConnInterface) GalleryClient {
	return &galleryClient{cc}
}

func (c *galleryClient) GetPhotos(ctx context.Context, in *PhotoFilter, opts ...grpc.CallOption) (Gallery_GetPhotosClient, error) {
	stream, err := c.cc.NewStream(ctx, &Gallery_ServiceDesc.Streams[0], "/gallery.Gallery/GetPhotos", opts...)
	if err != nil {
		return nil, err
	}
	x := &galleryGetPhotosClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Gallery_GetPhotosClient interface {
	Recv() (*Photo, error)
	grpc.ClientStream
}

type galleryGetPhotosClient struct {
	grpc.ClientStream
}

func (x *galleryGetPhotosClient) Recv() (*Photo, error) {
	m := new(Photo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GalleryServer is the server API for Gallery service.
// All implementations must embed UnimplementedGalleryServer
// for forward compatibility
type GalleryServer interface {
	// get all photos
	GetPhotos(*PhotoFilter, Gallery_GetPhotosServer) error
	mustEmbedUnimplementedGalleryServer()
}

// UnimplementedGalleryServer must be embedded to have forward compatible implementations.
type UnimplementedGalleryServer struct {
}

func (UnimplementedGalleryServer) GetPhotos(*PhotoFilter, Gallery_GetPhotosServer) error {
	return status.Errorf(codes.Unimplemented, "method GetPhotos not implemented")
}
func (UnimplementedGalleryServer) mustEmbedUnimplementedGalleryServer() {}

// UnsafeGalleryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GalleryServer will
// result in compilation errors.
type UnsafeGalleryServer interface {
	mustEmbedUnimplementedGalleryServer()
}

func RegisterGalleryServer(s grpc.ServiceRegistrar, srv GalleryServer) {
	s.RegisterService(&Gallery_ServiceDesc, srv)
}

func _Gallery_GetPhotos_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PhotoFilter)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GalleryServer).GetPhotos(m, &galleryGetPhotosServer{stream})
}

type Gallery_GetPhotosServer interface {
	Send(*Photo) error
	grpc.ServerStream
}

type galleryGetPhotosServer struct {
	grpc.ServerStream
}

func (x *galleryGetPhotosServer) Send(m *Photo) error {
	return x.ServerStream.SendMsg(m)
}

// Gallery_ServiceDesc is the grpc.ServiceDesc for Gallery service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gallery_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gallery.Gallery",
	HandlerType: (*GalleryServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetPhotos",
			Handler:       _Gallery_GetPhotos_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/gallery.proto",
}