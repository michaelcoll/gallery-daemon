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
)

const insertPhotoStmt = `
	INSERT INTO photos (hash, path) VALUES (?, ?);
	`

const selectByHashStmt = `
	SELECT hash,
		   path
	FROM photos
	WHERE hash = ?
	`

// insertPhotoIntoDB inserts a new line into the photos table
func insertPhotoIntoDB(db *sql.DB, photo *domain.Photo) error {
	stmt, err := db.Prepare(insertPhotoStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(photo.Hash, photo.Path)
	if err != nil {
		return err
	}

	return nil
}

// findByHash returns a photo having the given hash
func findByHash(db *sql.DB, hash string) (*domain.Photo, error) {
	stmt, err := db.Prepare(selectByHashStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var photo domain.Photo
	err = stmt.QueryRow(hash).Scan(&photo.Hash, &photo.Path)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &photo, nil
}
