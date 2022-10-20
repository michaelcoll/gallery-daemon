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

package service

import (
	"context"
	"log"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
)

type PhotoService struct {
	photoPath *string

	r repository.PhotoRepository
}

func New(r repository.PhotoRepository) PhotoService {
	return PhotoService{r: r}
}

func (s *PhotoService) indexImage(ctx context.Context, imagePath string) {
	photo := &model.Photo{Path: imagePath}
	extractData(photo)

	if err := s.r.CreateOrReplace(ctx, *photo); err != nil {
		log.Fatalf("Can't insert photo located at '%s' into database (%v)\n", imagePath, err)
	}
}

func (s *PhotoService) deleteImage(ctx context.Context, imagePath string) {
	if err := s.r.Delete(ctx, imagePath); err != nil {
		log.Fatalf("Can't delete photo with path '%s' (%v)\n", imagePath, err)
	}
}

func (s *PhotoService) CloseDb() {
	s.r.Close()
}
