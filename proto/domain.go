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

package gallery

import (
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
)

func (x *Photo) New(photo *model.Photo) *Photo {
	return &Photo{
		Hash: photo.Hash,
		Path: photo.Path,

		//DateTime:     toDateTime(photo.DateTime),
		Iso: int32(photo.Iso),
		//ExposureTime: toRational(photo.ExposureTime),
		XDimension: int32(photo.XDimension),
		YDimension: int32(photo.YDimension),
		Model:      photo.Model,
		//FocalLength:  toRational(photo.FocalLength),
	}
}
