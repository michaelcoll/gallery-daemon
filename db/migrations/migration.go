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

package migrations

import (
	"context"
	"database/sql"
	"log"
)

const testVersionTableExists = `
SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name='db_version';
`

const initSql = `
CREATE TABLE db_version
(
    version_number INTEGER NOT NULL
);
INSERT INTO db_version (version_number) VALUES (0);
`

const selectVersionStmt = `
SELECT version_number FROM db_version;
`

const updateVersionStmt = `
UPDATE db_version
SET version_number = ?
WHERE 1;
`

const v1Init = `
CREATE TABLE photos
(
    hash          TEXT PRIMARY KEY,
    path          TEXT NOT NULL,
    date_time     TEXT,
    iso           INTEGER,
    exposure_time TEXT,
    x_dimension   INTEGER,
    y_dimension   INTEGER,
    model         TEXT,
    focal_length  TEXT
);
`

var migrations = map[int]string{
	1: v1Init,
}

type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
}

// New creates a new instance of Migrations struct
func New(db DB) *Migrations {
	return &Migrations{db: db}
}

type Migrations struct {
	db DB
}

// Migrate migrates the database using the migration scripts provided
func (m *Migrations) Migrate(ctx context.Context) {
	initialized, err := m.isInitialized(ctx)
	if err != nil {
		log.Fatalf("Can't detect if database is initialized %v", err)
	}
	if initialized {
		version, err := m.getVersion(ctx)
		if err != nil {
			log.Fatalf("Can't read database version %v", err)
		}
		m.applyMigration(ctx, version)
	} else {
		m.createDBVersionTable(ctx)
		m.applyMigration(ctx, 0)
	}
}

// isInitialized checks if the table db_version is present in the current database
func (m *Migrations) isInitialized(ctx context.Context) (bool, error) {
	stmt, err := m.db.PrepareContext(ctx, testVersionTableExists)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var tablePresent int
	err = stmt.QueryRow().Scan(&tablePresent)
	if err != nil {
		return false, err
	}

	return tablePresent == 1, nil
}

// getVersion returns the current version of the schema
func (m *Migrations) getVersion(ctx context.Context) (int, error) {
	stmt, err := m.db.PrepareContext(ctx, selectVersionStmt)
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

// applyMigration a migration
func (m *Migrations) createDBVersionTable(ctx context.Context) {
	_, err := m.db.ExecContext(ctx, initSql)
	if err != nil {
		log.Fatalf("Could not create db_version table %v", err)
	}
}

// applyMigration a migration
func (m *Migrations) applyMigration(ctx context.Context, fromVersion int) {
	updStmt, err := m.db.PrepareContext(ctx, updateVersionStmt)
	if err != nil {
		log.Fatalf("Could not prepare Stmt : %v", err)
	}
	defer updStmt.Close()

	for version, script := range migrations {
		if version > fromVersion {
			_, err := m.db.ExecContext(ctx, script)
			if err != nil {
				log.Fatalf("Could not apply migration : %s, %v", script, err)
			}

			_, err = updStmt.ExecContext(ctx, version)
			if err != nil {
				log.Fatalf("Could not update version : %v", err)
			}
		}
	}
}
