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
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cozy/goexif2/exif"
	"github.com/cozy/goexif2/tiff"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
)

// extractExif extracts the EXIF data of a photo
func extractExif(photo *model.Photo) error {
	f, err := os.Open(photo.Path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Printf("Error reading EXIF data in file : %s (%v)\n", photo.Path, err)
	} else if err == nil {
		err := x.Walk(&walker{p: photo})
		if err != nil {
			panic(err)
		}
	}

	return nil
}

type walker struct {
	p *model.Photo
}

func (w *walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	if name == "DateTime" {
		dateTimeStr, _ := tag.StringVal()
		date, _ := time.Parse("2006:01:02 15:04:05", dateTimeStr)

		w.p.DateTime = date.Format("2006-01-02T15:04:05")
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
		value, _ := tag.StringVal()
		w.p.Model = strings.Trim(value, " ")
	} else if name == "FNumber" {
		if ratValue, err := tag.Rat(0); err == nil {
			floatValue, _ := ratValue.Float32()
			w.p.FNumber = fmt.Sprintf("f/%s", strconv.FormatFloat(float64(floatValue), 'f', -1, 32))
		}
	} else if name == "Orientation" {
		w.p.Orientation, _ = tag.Int(0)
	}
	return nil
}

func toString(rat *big.Rat) string {
	return fmt.Sprintf("%d/%d", rat.Num(), rat.Denom())
}

func fixPhotoAttributes(photo *model.Photo) {
	if photo.XDimension == 0 || photo.YDimension == 0 {
		reader, err := os.Open(photo.Path)
		if err != nil {
			log.Printf("Can't open file %s : %v\n", photo.Path, err)
		}
		defer reader.Close()

		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			log.Printf("Can't read image %s : %v\n", photo.Path, err)
		}

		photo.XDimension = im.Width
		photo.YDimension = im.Height
	}
	if photo.DateTime == "" {
		file, err := os.Stat(photo.Path)
		if err != nil {
			log.Printf("Can't open file %s : %v\n", photo.Path, err)
		}

		photo.DateTime = file.ModTime().Format("2006-01-02T15:04:05")
	}
	if photo.Orientation == 6 || photo.Orientation == 8 {
		x, y := photo.XDimension, photo.YDimension

		photo.XDimension = y
		photo.YDimension = x
	}
}
