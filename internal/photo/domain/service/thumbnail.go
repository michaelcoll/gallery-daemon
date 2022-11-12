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
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/webp"
)

const (
	quality = 80
)

func webpEncoder(photo *model.Photo) ([]byte, error) {
	var buf bytes.Buffer
	var img image.Image

	img, err := readRawImage(photo.Path)
	if err != nil {
		return nil, err
	}

	// Resize
	width, height := 0, 200
	if photo.Orientation == 6 || photo.Orientation == 8 {
		width, height = 200, 0
	}
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	// Rotate if necessary
	rotatedImg := rotate(photo, resizedImg)

	// Encode
	if err = webp.Encode(&buf, rotatedImg, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func readRawImage(imgPath string) (img image.Image, err error) {
	data, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}

	imgExtension := strings.ToLower(path.Ext(imgPath))
	if strings.Contains(imgExtension, "jpeg") || strings.Contains(imgExtension, "jpg") {
		img, err = jpeg.Decode(bytes.NewReader(data))
	}

	if err != nil || img == nil {
		errinfo := fmt.Sprintf("image file %s is corrupted: %v", imgPath, err)
		return nil, errors.New(errinfo)
	}

	return img, nil
}

func rotate(photo *model.Photo, image image.Image) image.Image {
	if photo.Orientation == 6 {
		return imaging.Rotate270(image)
	} else if photo.Orientation == 8 {
		return imaging.Rotate90(image)
	} else if photo.Orientation != 1 {
		fmt.Printf("Unsuported Orientation : %s\n", photo.Path)
	}

	return image
}
