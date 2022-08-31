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
	"database/sql"
	"github.com/michaelcoll/gallery-daemon/domain"
	"log"
)

// Index indexes a photo
func Index(db *sql.DB, photo *domain.Photo) error {
	_ = extractExif(photo)

	if err := insertPhotoIntoDB(db, photo); err != nil {
		log.Printf("\nhash %s\n", photo.Hash)
		return err
	}

	return nil
}
