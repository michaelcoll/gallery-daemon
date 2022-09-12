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

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/michaelcoll/gallery-daemon/db/connect"
	"github.com/michaelcoll/gallery-daemon/db/migrations"
	"github.com/michaelcoll/gallery-daemon/scanner"
	"github.com/michaelcoll/gallery-daemon/server"
)

var path = flag.String("f", ".", "The folder where the photos are.")

func main() {

	flag.Parse()

	fmt.Printf("Monitoring folder %s \n", color.GreenString(*path))

	ctx := context.Background()
	migrateDB(ctx)

	scanner.New().Scan(ctx, *path)

	server.Serve()

}

func migrateDB(ctx context.Context) {
	c := connect.Connect(false)
	defer c.Close()

	migrations.New(c).Migrate(ctx)
}
