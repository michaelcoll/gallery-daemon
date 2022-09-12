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
	"fmt"
	"github.com/cozy/goexif2/exif"
	"github.com/cozy/goexif2/tiff"
	"github.com/michaelcoll/gallery-daemon/db"
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
	ExposureTime string
	XDimension   int
	YDimension   int
	Model        string
	FocalLength  string
}

func (data *Photo) Walk(name exif.FieldName, tag *tiff.Tag) error {
	if name == "DateTime" {
		dateTimeStr, _ := tag.StringVal()

		data.DateTime, _ = time.Parse("2006:01:02 15:04:05", dateTimeStr)
	} else if name == "ISOSpeedRatings" {
		data.Iso, _ = tag.Int(0)
	} else if name == "ExposureTime" {
		value, _ := tag.Rat(0)
		data.ExposureTime = toString(value)
	} else if name == "PixelXDimension" {
		data.XDimension, _ = tag.Int(0)
	} else if name == "PixelYDimension" {
		data.YDimension, _ = tag.Int(0)
	} else if name == "Model" {
		data.Model, _ = tag.StringVal()
	} else if name == "FocalLength" {
		value, _ := tag.Rat(0)
		data.FocalLength = toString(value)
	}
	return nil
}

func toString(rat *big.Rat) string {
	return fmt.Sprintf("%d/%d", rat.Num(), rat.Denom())
}

func (data *Photo) ToDBInsert() (*db.CreatePhotoParams, error) {
	params := db.CreatePhotoParams{
		Hash: data.Hash,
		Path: data.Path,
	}

	err := params.DateTime.Scan(data.DateTime)
	if err != nil {
		return nil, err
	}
	err = params.Iso.Scan(data.Iso)
	if err != nil {
		return nil, err
	}
	err = params.ExposureTime.Scan(data.ExposureTime)
	if err != nil {
		return nil, err
	}
	err = params.XDimension.Scan(data.XDimension)
	if err != nil {
		return nil, err
	}
	err = params.YDimension.Scan(data.YDimension)
	if err != nil {
		return nil, err
	}
	err = params.Model.Scan(data.Model)
	if err != nil {
		return nil, err
	}
	err = params.FocalLength.Scan(data.FocalLength)
	if err != nil {
		return nil, err
	}

	return &params, nil
}
