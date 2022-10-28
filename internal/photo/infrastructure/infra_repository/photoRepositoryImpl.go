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

package infra_repository

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/consts"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/db"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/sqlc"
)

const BufferSize = 1024 * 1024 * 2
const webPContentType = "image/webp"

type PhotoDBRepository struct {
	repository.PhotoRepository

	databaseLocation string

	c *sql.DB
	q *sqlc.Queries
}

func New(localDb bool, photosPath string) *PhotoDBRepository {
	var databaseLocation string
	if localDb {
		databaseLocation = "."
	} else {
		databaseLocation = photosPath
	}

	connection := db.Connect(false, databaseLocation)
	db.New(connection).Migrate()

	return &PhotoDBRepository{databaseLocation: databaseLocation, q: sqlc.New(), c: connection}
}

func (r *PhotoDBRepository) Connect(readOnly bool) {
	r.c = db.Connect(readOnly, r.databaseLocation)
}

func (r *PhotoDBRepository) Close() {
	r.c.Close()
}

func (r *PhotoDBRepository) CreateOrReplace(ctx context.Context, photo model.Photo) error {
	params, err := r.toInfra(photo)
	if err != nil {
		return err
	}

	if err := r.q.CreateOrReplacePhoto(ctx, r.c, params); err != nil {
		return err
	}

	return nil
}

func (r *PhotoDBRepository) Get(ctx context.Context, hash string) (model.Photo, error) {
	photo, err := r.q.GetPhoto(ctx, r.c, hash)
	if err == sql.ErrNoRows {
		return model.Photo{}, status.Error(codes.NotFound, "media not found")
	}
	if err != nil {
		return model.Photo{}, err
	}
	domain, err := r.toDomainGet(photo)
	if err != nil {
		return model.Photo{}, err
	}

	return domain, nil
}

func (r *PhotoDBRepository) Exists(ctx context.Context, hash string) bool {
	count, err := r.q.CountPhotoByHash(ctx, r.c, hash)
	if err != nil {
		return false
	}

	return count == 1
}

func (r *PhotoDBRepository) List(ctx context.Context, page int32, pageSize int32) ([]model.Photo, error) {
	list, err := r.q.List(ctx, r.c, sqlc.ListParams{
		Limit:  int64(pageSize),
		Offset: int64(page * pageSize),
	})
	if err != nil {
		return nil, err
	}

	photos := make([]model.Photo, len(list))
	for i, photo := range list {
		domain, err := r.toDomainList(photo)
		if err != nil {
			return nil, err
		}
		photos[i] = domain
	}

	return photos, nil
}

func (r *PhotoDBRepository) ReadContent(ctx context.Context, hash string, reader repository.ImageReader) error {
	photo, err := r.q.GetPhoto(ctx, r.c, hash)
	if err == sql.ErrNoRows {
		return status.Error(codes.NotFound, "media not found")
	}
	if err != nil {
		return err
	}

	contentType, err := detectContentType(photo.Path)
	if err != nil {
		return err
	}

	f, err := os.Open(fmt.Sprintf("%s%s", r.databaseLocation, photo.Path))
	if err != nil {
		return err
	}
	defer f.Close()

	fReader := bufio.NewReader(f)
	buf := make([]byte, BufferSize)

	for {
		n, err := fReader.Read(buf)

		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		err = reader.ReadChunk(buf[0:n], contentType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PhotoDBRepository) ReadThumbnail(ctx context.Context, hash string, reader repository.ImageReader) error {
	thumbnailBytes, err := r.q.GetThumbnail(ctx, r.c, hash)
	if err == sql.ErrNoRows {
		return status.Error(codes.NotFound, "media not found")
	}
	if err != nil {
		return err
	}

	if err = reader.ReadChunk(thumbnailBytes, webPContentType); err != nil {
		return err
	}

	return nil
}

func (r *PhotoDBRepository) SetThumbnail(ctx context.Context, hash string, thumbnail []byte) error {
	if err := r.q.UpdateThumbnail(ctx, r.c, sqlc.UpdateThumbnailParams{
		Thumbnail: thumbnail,
		Hash:      hash,
	}); err != nil {
		return err
	}

	return nil
}

func detectContentType(photoPath string) (string, error) {
	for ext, contentType := range consts.ExtensionsAndContentTypesMap {
		if strings.HasSuffix(photoPath, ext) {
			return contentType, nil
		}
	}

	return "", errors.New("content type not supported")
}

func (r *PhotoDBRepository) Delete(ctx context.Context, path string) error {
	err := r.q.DeletePhotoByPath(ctx, r.c, path)
	if err != nil {
		return err
	}

	return nil
}

func (r *PhotoDBRepository) DeleteAllPhotoInPath(ctx context.Context, path string) error {
	err := r.q.DeleteAllPhotoInPath(ctx, r.c, fmt.Sprintf("'%s%%'", strings.ReplaceAll(path, r.databaseLocation, "")))
	if err != nil {
		return err
	}

	return nil
}

func (r *PhotoDBRepository) DeleteAll(ctx context.Context) error {
	err := r.q.DeleteAllPhotos(ctx, r.c)
	if err != nil {
		return err
	}

	return nil
}

func (r *PhotoDBRepository) CountPhotos(ctx context.Context) (int, error) {
	count, err := r.q.CountPhotos(ctx, r.c)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
