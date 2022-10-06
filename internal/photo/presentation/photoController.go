/*
 * Copyright (c) 2022 Michaël COLL.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package presentation

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"google.golang.org/grpc"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
	pb "github.com/michaelcoll/gallery-daemon/proto"
)

const port = ":9000"

type PhotoController struct {
	pb.UnimplementedGalleryServer

	r repository.PhotoRepository
}

func New(r repository.PhotoRepository) PhotoController {
	return PhotoController{r: r}
}

func (c *PhotoController) Serve() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterGalleryServer(grpcServer, c)

	fmt.Printf("Listening on port %s\n", color.GreenString(port))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// GetPhotos returns all photos by given filter
func (c *PhotoController) GetPhotos(ctx context.Context, filter *pb.ListFilter) (*pb.PhotosResponse, error) {

	c.r.Connect(true)
	defer c.r.Close()

	list, err := c.r.List(ctx)
	if err != nil {
		return nil, err
	}

	responseList := make([]*pb.Photo, len(list))
	for i, photo := range list {
		responseList[i] = toGrpc(photo)
	}

	return &pb.PhotosResponse{Photos: responseList}, nil
}

func (c *PhotoController) GetByHash(ctx context.Context, filter *pb.HashFilter) (*pb.Photo, error) {

	c.r.Connect(true)
	defer c.r.Close()

	photo, err := c.r.Get(ctx, filter.Hash)
	if err != nil {
		return nil, err
	}

	return toGrpc(photo), nil
}

func (c *PhotoController) ContentByHash(filter *pb.HashFilter, stream pb.Gallery_ContentByHashServer) error {

	c.r.Connect(true)
	defer c.r.Close()

	err := c.r.ReadContent(stream.Context(), filter.Hash, streamReader{stream: stream})
	if err != nil {
		return err
	}

	return nil
}

type streamReader struct {
	stream pb.Gallery_ContentByHashServer
}

func (r streamReader) ReadChunk(bytes []byte) error {
	err := r.stream.Send(&pb.PhotoChunk{
		Data: bytes,
	})
	if err != nil {
		return err
	}

	return nil
}
