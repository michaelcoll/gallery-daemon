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
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/sqlc"
)

func (r *PhotoDBRepository) toInfra(photo model.Photo) (sqlc.CreatePhotoParams, error) {
	params := sqlc.CreatePhotoParams{
		Hash: photo.Hash,
		Path: photo.Path,
	}

	if err := params.DateTime.Scan(photo.DateTime); err != nil {
		return sqlc.CreatePhotoParams{}, err
	}
	if photo.Iso != 0 {
		if err := params.Iso.Scan(photo.Iso); err != nil {
			return sqlc.CreatePhotoParams{}, err
		}
	}
	if err := params.ExposureTime.Scan(photo.ExposureTime); err != nil {
		return sqlc.CreatePhotoParams{}, err
	}
	if err := params.XDimension.Scan(photo.XDimension); err != nil {
		return sqlc.CreatePhotoParams{}, err
	}
	if err := params.YDimension.Scan(photo.YDimension); err != nil {
		return sqlc.CreatePhotoParams{}, err
	}
	if err := params.Model.Scan(photo.Model); err != nil {
		return sqlc.CreatePhotoParams{}, err
	}
	if err := params.FNumber.Scan(photo.FNumber); err != nil {
		return sqlc.CreatePhotoParams{}, err
	}

	return params, nil
}

func (r *PhotoDBRepository) toDomain(photo sqlc.Photo) (model.Photo, error) {

	m := &model.Photo{
		Hash: photo.Hash,
		Path: photo.Path,
	}

	if photo.DateTime.Valid {
		m.DateTime = photo.DateTime.String
	}
	if photo.Iso.Valid {
		m.Iso = int(photo.Iso.Int32)
	}
	if photo.ExposureTime.Valid {
		m.ExposureTime = photo.ExposureTime.String
	}
	if photo.XDimension.Valid {
		m.XDimension = int(photo.XDimension.Int32)
	}
	if photo.YDimension.Valid {
		m.YDimension = int(photo.YDimension.Int32)
	}
	if photo.Model.Valid {
		m.Model = photo.Model.String
	}
	if photo.FNumber.Valid {
		m.FNumber = photo.FNumber.String
	}

	return *m, nil
}
