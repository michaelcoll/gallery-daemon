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
	"fmt"
	"github.com/michaelcoll/gallery-daemon/database"
	"github.com/michaelcoll/gallery-daemon/domain"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func Scan(path string) error {
	images := getImageFiles(path, false)

	db := database.Connect(true)

	fmt.Printf("Syncing database ... ")

	var imagesToInsert []domain.Photo
	for _, scannedImage := range images {
		photo, err := findByHash(db, scannedImage.Hash)
		if err != nil {
			return err
		}

		if photo == nil {
			imagesToInsert = append(imagesToInsert, scannedImage)
		}
	}
	db.Close()

	if len(imagesToInsert) > 0 {
		fmt.Printf("Inserting %d photo(s) ... ", len(imagesToInsert))

		db = database.Connect(false)
		defer db.Close()

		for _, photo := range imagesToInsert {
			if err := Index(db, &photo); err != nil {
				return err
			}
		}
	}

	fmt.Println("Done.")

	return nil
}

func getImageFiles(path string, isSubFolder bool) []domain.Photo {
	var images []domain.Photo

	if !isSubFolder {
		fmt.Printf("Finding all images in folder %s ... ", path)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			for _, image := range getImageFiles(filepath.Join(path, file.Name()), true) {
				images = append(images, image)
			}
		} else if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".JPG") || strings.HasSuffix(file.Name(), ".JPEG") {
			imagePath := filepath.Join(path, file.Name())
			hash, err := sha(imagePath)
			if err != nil {
				log.Printf("\nCan't calculate hash for file : %s", path)
			}
			images = append(images, domain.Photo{Hash: hash, Path: imagePath})
		}
	}

	if !isSubFolder {
		fmt.Printf("Found %d image(s).\n", len(images))
	}

	return images
}
