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

package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(readOnly bool) *sql.DB {
	db, err := sql.Open("sqlite3", getDBUrl(readOnly))
	if err != nil {
		log.Fatalf("Can't open database %v", err)
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
