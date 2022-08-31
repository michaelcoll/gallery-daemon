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

package domain

import (
	"github.com/cozy/goexif2/exif"
	"github.com/cozy/goexif2/tiff"
	"math/big"
	"time"
)

type Photo struct {

	// Main

	Hash string
	Path string

	// EXIF

	DateTime     time.Time
	Iso          int
	ExposureTime *big.Rat
	XDimension   int
	YDimension   int
	Model        string
	FocalLength  *big.Rat
}

func (data *Photo) Walk(name exif.FieldName, tag *tiff.Tag) error {
	if name == "DateTime" {
		dateTimeStr, _ := tag.StringVal()

		data.DateTime, _ = time.Parse("2006:01:02 15:04:05", dateTimeStr)
	} else if name == "ISOSpeedRatings" {
		data.Iso, _ = tag.Int(0)
	} else if name == "ExposureTime" {
		data.ExposureTime, _ = tag.Rat(0)
	} else if name == "PixelXDimension" {
		data.XDimension, _ = tag.Int(0)
	} else if name == "PixelYDimension" {
		data.YDimension, _ = tag.Int(0)
	} else if name == "Model" {
		data.Model, _ = tag.StringVal()
	} else if name == "FocalLength" {
		data.FocalLength, _ = tag.Rat(0)
	}
	return nil
}
