// Code generated by sqlc-addon. DO NOT EDIT.
// versions:
//   sqlc-addon v1.3.0

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(readOnly bool, baseLocation string) *sql.DB {
	db, err := sql.Open("sqlite3", getDBUrl(readOnly, baseLocation))
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}

	return db
}

func getDBUrl(readOnly bool, baseLocation string) string {

	var options string
	if readOnly {
		options = "cache=shared&mode=ro"
	} else {
		options = "cache=shared&mode=rwc&_auto_vacuum=full&_journal_mode=WAL"
	}

	return fmt.Sprintf("file:%s/%s?%s", baseLocation, "photos.db", options)
}
