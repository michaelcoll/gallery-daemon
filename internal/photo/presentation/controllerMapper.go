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
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"

	photov1 "github.com/michaelcoll/gallery-proto/gen/proto/go/photo/v1"
)

func toGrpc(photo model.Photo) *photov1.Photo {
	return &photov1.Photo{
		Hash: photo.Hash,
		Path: photo.Path,

		DateTime:     photo.DateTime,
		Iso:          uint32(photo.Iso),
		ExposureTime: photo.ExposureTime,
		XDimension:   uint32(photo.XDimension),
		YDimension:   uint32(photo.YDimension),
		Model:        photo.Model,
		FNumber:      photo.FNumber,
	}
}
