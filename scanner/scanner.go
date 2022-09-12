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

package scanner

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/michaelcoll/gallery-daemon/db"
	"github.com/michaelcoll/gallery-daemon/db/connect"
	"github.com/michaelcoll/gallery-daemon/domain"
)

type Scanner struct {
	q *db.Queries
}

func New() *Scanner {
	c := connect.Connect(false)

	return &Scanner{q: db.New(c)}
}

func (s *Scanner) Scan(ctx context.Context, path string) {
	var wg sync.WaitGroup
	images := getImageFiles(path, false, &wg)
	wg.Wait()

	fmt.Printf("Syncing database ... ")

	var imagesToInsert []*domain.Photo
	for _, scannedImage := range images {
		_, err := s.q.GetPhoto(ctx, scannedImage.Hash)
		if err == sql.ErrNoRows {
			imagesToInsert = append(imagesToInsert, scannedImage)
		} else if err != nil {
			log.Fatalf("Can't find photo in database, %v", err)
		} else {
			log.Printf("Image %s already present in database.", scannedImage.Path)
		}
	}

	if len(imagesToInsert) > 0 {
		fmt.Printf("Inserting %d photo(s) ... \n", len(imagesToInsert))

		bar := progressbar.Default(int64(len(imagesToInsert)))
		for _, photo := range imagesToInsert {
			params, err := photo.ToDBInsert()
			if err != nil {
				log.Fatalf("Can't convert photo data, %v", err)
			}

			err = s.q.CreatePhoto(ctx, *params)
			if err != nil {
				log.Fatalf("Can't insert photo into database (%v), %v", err, params)
			}
			_ = bar.Add(1)
		}
	}
}

func getImageFiles(path string, isSubFolder bool, wg *sync.WaitGroup) []*domain.Photo {
	var images []*domain.Photo

	if !isSubFolder {
		fmt.Printf("Finding all images in folder %s ... ", path)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			for _, image := range getImageFiles(filepath.Join(path, file.Name()), true, wg) {
				images = append(images, image)
			}
		} else if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".JPG") || strings.HasSuffix(file.Name(), ".JPEG") {
			imagePath := filepath.Join(path, file.Name())
			photo := &domain.Photo{Path: imagePath}

			wg.Add(1)
			go func() {
				extractData(photo)
				wg.Done()
			}()

			images = append(images, photo)
		}
	}

	if !isSubFolder {
		fmt.Printf("Found %d image(s).\n", len(images))
	}

	return images
}

func extractData(photo *domain.Photo) {
	hash, err := sha(photo.Path)
	if err != nil {
		log.Printf("\nCan't calculate hash for file : %s", photo.Path)
	}

	photo.Hash = hash

	_ = extractExif(photo)
}
