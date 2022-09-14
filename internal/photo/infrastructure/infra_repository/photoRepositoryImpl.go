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
	"context"
	"database/sql"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/db"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/sqlc"
)

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

func (r *PhotoDBRepository) Create(ctx context.Context, photo model.Photo) error {
	params, err := r.toInfra(photo)
	if err != nil {
		return err
	}

	if err := r.q.CreatePhoto(ctx, r.c, params); err != nil {
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
