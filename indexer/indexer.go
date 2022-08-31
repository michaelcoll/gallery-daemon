package indexer

import (
	"database/sql"
	"github.com/michaelcoll/gallery-daemon/domain"
	"log"
)

// Index indexes a photo
func Index(db *sql.DB, photo *domain.Photo) error {
	_ = extractExif(photo)

	if err := insertPhotoIntoDB(db, photo); err != nil {
		log.Printf("\nhash %s\n", photo.Hash)
		return err
	}

	return nil
}
