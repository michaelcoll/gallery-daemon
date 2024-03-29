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
	domain, err := r.toDomain(photo)
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

func (r *PhotoDBRepository) List(ctx context.Context, offset uint32, limit uint32) ([]model.Photo, error) {
	list, err := r.q.List(ctx, r.c, sqlc.ListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, err
	}

	photos := make([]model.Photo, len(list))
	for i, photo := range list {
		domain, err := r.toDomain(photo)
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

func (r *PhotoDBRepository) ReadThumbnail(ctx context.Context, hash string, width uint32, height uint32, reader repository.ImageReader) error {

	var w, h uint32
	if width > 0 && height > 0 {
		w = width
	} else if width == 0 && height == 0 {
		w = 200
	} else {
		w, h = width, height
	}

	thumbnailBytes, err := r.q.GetThumbnail(ctx, r.c, sqlc.GetThumbnailParams{
		Hash:   hash,
		Width:  int64(w),
		Height: int64(h),
	})
	if err == sql.ErrNoRows {
		thumbnailBytes, err = r.createAndUpdateThumbnail(ctx, hash, w, h)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	if err = reader.ReadChunk(thumbnailBytes, webPContentType); err != nil {
		return err
	}

	return nil
}

func (r *PhotoDBRepository) createAndUpdateThumbnail(ctx context.Context, hash string, width uint32, height uint32) ([]byte, error) {
	photo, err := r.q.GetPhoto(ctx, r.c, hash)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "media not found")
	}

	path := fmt.Sprintf("%s%s", r.databaseLocation, photo.Path)
	orientation := getOrientation(photo)

	if thumbnail, err := webpEncoder(path, orientation); err != nil {
		return nil, status.Errorf(codes.Internal, "Error while creating the thumbnail of the file %s : %v\n", photo.Path, err)
	} else {
		if err := r.createOrReplaceThumbnail(ctx, photo.Hash, width, height, thumbnail); err != nil {
			return nil, status.Errorf(codes.Internal, "Error save thumbnail in database (%v).\n", err)
		}

		return thumbnail, nil
	}
}

func getOrientation(photo sqlc.Photo) uint {
	if photo.Orientation.Valid {
		return uint(photo.Orientation.Int64)
	}
	return 1
}

func (r *PhotoDBRepository) createOrReplaceThumbnail(ctx context.Context, hash string, width uint32, height uint32, thumbnail []byte) error {
	if err := r.q.CreateOrReplaceThumbnail(ctx, r.c, sqlc.CreateOrReplaceThumbnailParams{
		Hash:      hash,
		Width:     int64(width),
		Height:    int64(height),
		Thumbnail: thumbnail,
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

func (r *PhotoDBRepository) CountPhotos(ctx context.Context) (uint32, error) {
	count, err := r.q.CountPhotos(ctx, r.c)
	if err != nil {
		return 0, err
	}

	return uint32(count), nil
}
