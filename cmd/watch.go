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

package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/michaelcoll/gallery-daemon/internal/photo"
)

// indexCmd represents the index command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch the given folder for updates",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Watching folder %s \n", color.GreenString(folder))

		photo.NewForIndex(localDb, folder).GetPhotoService().Watch(folder)
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
