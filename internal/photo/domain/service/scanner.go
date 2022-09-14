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

package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/cozy/goexif2/exif"
	"github.com/cozy/goexif2/tiff"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
)

// Scan scans the folder given in parameter and fill the database with image info and EXIF data found on JPEGs
func (s *PhotoService) Scan(ctx context.Context, path string) {

	s.r.Connect(false)
	defer s.r.Close()

	bar := progressbar.Default(-1, fmt.Sprintf("Finding all images in folder %s ... ", path))
	var wg sync.WaitGroup
	images := getImageFiles(path, &wg, bar)
	wg.Wait()
	_ = bar.Clear()

	bar = progressbar.Default(-1, "Syncing database ... ")
	var imagesToInsert []*model.Photo
	for _, scannedImage := range images {
		_ = bar.Add(1)
		if !s.r.Exists(ctx, scannedImage.Hash) {
			imagesToInsert = append(imagesToInsert, scannedImage)
		}
	}
	_ = bar.Clear()

	if len(imagesToInsert) > 0 {
		bar := progressbar.Default(int64(len(imagesToInsert)), "Updating database ...")
		for _, photo := range imagesToInsert {
			err := s.r.Create(ctx, *photo)
			if err != nil {
				log.Fatalf("Can't insert photo into database (%v)", err)
			}
			_ = bar.Add(1)
		}
		_ = bar.Clear()
	}

}

func getImageFiles(path string, wg *sync.WaitGroup, bar *progressbar.ProgressBar) []*model.Photo {
	var images []*model.Photo

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			for _, image := range getImageFiles(filepath.Join(path, file.Name()), wg, bar) {
				images = append(images, image)
			}
		} else if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".JPG") || strings.HasSuffix(file.Name(), ".JPEG") {
			imagePath := filepath.Join(path, file.Name())
			photo := &model.Photo{Path: imagePath}

			wg.Add(1)
			go func() {
				extractData(photo)
				_ = bar.Add(1)
				wg.Done()
			}()

			images = append(images, photo)
		}
	}

	return images
}

func extractData(photo *model.Photo) {
	hash, err := sha(photo.Path)
	if err != nil {
		log.Printf("\nCan't calculate hash for file : %s", photo.Path)
	}

	photo.Hash = hash

	_ = extractExif(photo)
}

// sha calculate the SHA256 of a file
func sha(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// extractExif extracts the EXIF data of a photo
func extractExif(photo *model.Photo) error {
	f, err := os.Open(photo.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return err
	} else {
		err := x.Walk(&walker{p: photo})
		if err != nil {
			return err
		}
		return nil
	}
}

type walker struct {
	p *model.Photo
}

func (w *walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	if name == "DateTime" {
		w.p.DateTime, _ = tag.StringVal()
	} else if name == "ISOSpeedRatings" {
		w.p.Iso, _ = tag.Int(0)
	} else if name == "ExposureTime" {
		if value, err := tag.Rat(0); err == nil {
			w.p.ExposureTime = toString(value)
		}
	} else if name == "PixelXDimension" {
		w.p.XDimension, _ = tag.Int(0)
	} else if name == "PixelYDimension" {
		w.p.YDimension, _ = tag.Int(0)
	} else if name == "Model" {
		w.p.Model, _ = tag.StringVal()
	} else if name == "MaxApertureValue" {
		if value, err := tag.Rat(0); err == nil {
			w.p.Aperture, _ = value.Float32()
		}
	}
	return nil
}

func toString(rat *big.Rat) string {
	return fmt.Sprintf("%d/%d", rat.Num(), rat.Denom())
}
