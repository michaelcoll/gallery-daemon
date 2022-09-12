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

package server

import (
	"fmt"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/michaelcoll/gallery-daemon/proto"
)

const port = ":9000"

// server is used to implement customer.CustomerServer.
type server struct {
	pb.UnimplementedGalleryServer
}

// GetPhotos returns all photos by given filter
func (s *server) GetPhotos(filter *pb.PhotoFilter, stream pb.Gallery_GetPhotosServer) error {
	return nil
}

func Serve() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterGalleryServer(s, &server{})

	fmt.Printf("Listening on port %s \n", color.GreenString(port))
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
