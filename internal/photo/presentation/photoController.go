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
	"strconv"

	"github.com/fatih/color"
	"google.golang.org/grpc"

	photov1 "github.com/michaelcoll/gallery-proto/gen/proto/go/photo/v1"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
)

type PhotoController struct {
	photov1.UnimplementedPhotoServiceServer

	r    repository.PhotoRepository
	port uint32
}

func New(r repository.PhotoRepository, param model.ServeParameters) PhotoController {
	return PhotoController{r: r, port: param.GrpcPort}
}

func (c *PhotoController) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	grpcServer := grpc.NewServer()
	photov1.RegisterPhotoServiceServer(grpcServer, c)

	fmt.Printf("%s Listening on 0.0.0.0:%s\n", color.GreenString("✓"), color.GreenString(strconv.Itoa(int(c.port))))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// GetPhotos returns all photos by given filter
func (c *PhotoController) GetPhotos(ctx context.Context, request *photov1.GetPhotosRequest) (*photov1.GetPhotosResponse, error) {

	var pageSize uint32
	if request.PageSize == 0 {
		pageSize = 25
	} else {
		pageSize = request.PageSize
	}

	list, err := c.r.List(ctx, request.Page, pageSize)
	if err != nil {
		return nil, err
	}

	responseList := make([]*photov1.Photo, len(list))
	for i, photo := range list {
		responseList[i] = toGrpc(photo)
	}

	return &photov1.GetPhotosResponse{Photos: responseList}, nil
}

func (c *PhotoController) GetByHash(ctx context.Context, request *photov1.GetByHashRequest) (*photov1.GetByHashResponse, error) {

	photo, err := c.r.Get(ctx, request.Hash)
	if err != nil {
		return nil, err
	}

	return &photov1.GetByHashResponse{Photo: toGrpc(photo)}, nil
}

func (c *PhotoController) ExistsByHash(ctx context.Context, request *photov1.ExistsByHashRequest) (*photov1.ExistsByHashResponse, error) {
	if len(request.Hash) < 40 {
		return &photov1.ExistsByHashResponse{Exists: false}, nil
	}

	exists := c.r.Exists(ctx, request.Hash)

	return &photov1.ExistsByHashResponse{Exists: exists}, nil
}

func (c *PhotoController) ContentByHash(filter *photov1.ContentByHashRequest, stream photov1.PhotoService_ContentByHashServer) error {

	err := c.r.ReadContent(stream.Context(), filter.Hash, &streamReader{stream: stream})
	if err != nil {
		return err
	}

	return nil
}

func (c *PhotoController) ThumbnailByHash(req *photov1.ThumbnailByHashRequest, stream photov1.PhotoService_ThumbnailByHashServer) error {

	err := c.r.ReadThumbnail(stream.Context(), req.Hash, req.Width, req.Height, &streamReaderThumbnail{stream: stream})
	if err != nil {
		return err
	}

	return nil
}

type streamReader struct {
	stream photov1.PhotoService_ContentByHashServer
}

func (r *streamReader) ReadChunk(bytes []byte, contentType string) error {
	err := r.stream.Send(&photov1.PhotoServiceContentByHashResponse{
		Data:        bytes,
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}

type streamReaderThumbnail struct {
	stream photov1.PhotoService_ThumbnailByHashServer
}

func (r *streamReaderThumbnail) ReadChunk(bytes []byte, contentType string) error {
	err := r.stream.Send(&photov1.PhotoServiceThumbnailByHashResponse{
		Data:        bytes,
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}
