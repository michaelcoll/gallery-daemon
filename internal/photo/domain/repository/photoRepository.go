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

package repository

import (
	"context"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
)

type ImageReader interface {
	ReadChunk([]byte, string) error
}

type PhotoRepository interface {
	// Connect Opens a database connection
	Connect(readOnly bool)
	// Close Closes the database connection
	Close()

	CreateOrReplace(context.Context, model.Photo) error
	Get(ctx context.Context, hash string) (model.Photo, error)
	ReadContent(ctx context.Context, hash string, reader ImageReader) error
	ReadThumbnail(ctx context.Context, hash string, reader ImageReader) error
	SetThumbnail(ctx context.Context, hash string, thumbnail []byte) error
	Exists(ctx context.Context, hash string) bool
	List(ctx context.Context, page int32, pageSize int32) ([]model.Photo, error)
	Delete(ctx context.Context, path string) error
	DeleteAllPhotoInPath(ctx context.Context, path string) error
	DeleteAll(ctx context.Context) error
	CountPhotos(ctx context.Context) (int, error)
}
