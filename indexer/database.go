package indexer

import (
	"database/sql"
	"github.com/michaelcoll/gallery-daemon/domain"
)

const insertPhotoStmt = `
	INSERT INTO photos (hash, path) VALUES (?, ?);
	`

const selectByHashStmt = `
	SELECT hash,
		   path
	FROM photos
	WHERE hash = ?
	`

// insertPhotoIntoDB inserts a new line into the photos table
func insertPhotoIntoDB(db *sql.DB, photo *domain.Photo) error {
	stmt, err := db.Prepare(insertPhotoStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(photo.Hash, photo.Path)
	if err != nil {
		return err
	}

	return nil
}

// findByHash returns a photo having the given hash
func findByHash(db *sql.DB, hash string) (*domain.Photo, error) {
	stmt, err := db.Prepare(selectByHashStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var photo domain.Photo
	err = stmt.QueryRow(hash).Scan(&photo.Hash, &photo.Path)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &photo, nil
}
