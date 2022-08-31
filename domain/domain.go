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
