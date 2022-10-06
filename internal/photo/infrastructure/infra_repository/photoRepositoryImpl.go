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
	"io"
	"os"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/db"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/sqlc"
)

const BUFFER_SIZE = 1024 * 1024 * 2

type PhotoDBRepository struct {
	repository.PhotoRepository

	c *sql.DB
	q *sqlc.Queries
}

func New() *PhotoDBRepository {
	connection := db.Connect(false)
	defer connection.Close()
	db.New(connection).Migrate()

	return &PhotoDBRepository{q: sqlc.New()}
}

func (r *PhotoDBRepository) Connect(readOnly bool) {
	r.c = db.Connect(readOnly)
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

func (r *PhotoDBRepository) List(ctx context.Context) ([]model.Photo, error) {
	list, err := r.q.List(ctx, r.c)
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
	if err != nil {
		return err
	}

	f, err := os.Open(photo.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	fReader := bufio.NewReader(f)
	buf := make([]byte, BUFFER_SIZE)

	for {
		n, err := fReader.Read(buf)

		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		err = reader.ReadChunk(buf[0:n])
		if err != nil {
			return err
		}
	}

	return nil
}
