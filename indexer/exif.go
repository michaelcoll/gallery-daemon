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

package indexer

import (
	"github.com/cozy/goexif2/exif"
	"github.com/michaelcoll/gallery-daemon/domain"
	"os"
)

// extractExif extracts the EXIF data of a photo
func extractExif(photo *domain.Photo) error {
	f, err := os.Open(photo.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return err
	} else {
		err := x.Walk(photo)
		if err != nil {
			return err
		}
		return nil
	}
}
