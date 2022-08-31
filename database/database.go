package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const dbVersion = 1

const initSql = `
	CREATE TABLE db_version
	(
		version_number INTEGER
	);
	INSERT INTO db_version (version_number) VALUES (1);
	
	CREATE TABLE photos
	(
		hash TEXT PRIMARY KEY,
		path TEXT
	);`

const selectVersionStmt = `
	SELECT version_number FROM db_version;`

func Connect(readOnly bool) *sql.DB {
	db, err := sql.Open("sqlite3", getDBUrl(readOnly))
	if err != nil {
		log.Fatal(err)
	}

	version, err := getVersion(db)
	if err != nil && !readOnly {
		initDB(db)
		version = dbVersion
	}

	if version != dbVersion {
		log.Fatalf("DB version mismatch ! (current: %d, target: %d)\n", version, dbVersion)
	}

	return db
}

func getDBUrl(readOnly bool) string {
	if readOnly {
		return "file:./photos.db?cache=shared&mode=ro"
	} else {
		return "file:./photos.db?cache=shared&mode=rwc&_auto_vacuum=full"
	}
}

func getVersion(db *sql.DB) (int, error) {
	stmt, err := db.Prepare(selectVersionStmt)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var version int
	err = stmt.QueryRow().Scan(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func initDB(db *sql.DB) {
	_, err := db.Exec(initSql)
	if err != nil {
		log.Fatal(err)
	}
}
