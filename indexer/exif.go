package indexer

import (
	"github.com/cozy/goexif2/exif"
	"github.com/michaelcoll/gallery-daemon/domain"
	"os"
)

// extractExif extracts the EXIF data of a photo
func extractExif(photo *domain.Photo) error {
	f, err := os.Open(photo.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return err
	} else {
		err := x.Walk(photo)
		if err != nil {
			return err
		}
		return nil
	}
}
