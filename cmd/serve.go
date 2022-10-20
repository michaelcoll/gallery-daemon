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
	"context"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/banner"
	"github.com/spf13/cobra"

	"github.com/michaelcoll/gallery-daemon/internal/photo"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long: `
Starts the daemon in server mode.

In this mode it will :
 - index the images if the database is not up-to-date
 - register the daemon to the backend
 - watch for file changes
 - serve backend requests`,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print(rootCmd.Version, banner.Serve)

		module := photo.NewForServe(localDb, folder, model.ServeParameters{
			GrpcPort:      grpcPort,
			ExternalHost:  externalHost,
			DaemonName:    name,
			DaemonVersion: Version,
		})

		// Indexation
		photoService := module.GetPhotoService()
		if reIndex {
			photoService.ReIndex(context.Background(), folder)
		} else {
			photoService.Index(context.Background(), folder)
		}
		photoService.CloseDb()

		// Registration
		go module.GetRegisterService().Register()

		// Watch for file changes
		go photoService.Watch(folder)

		// Serving backend requests
		module.GetController().Serve()
	},
}

var grpcPort int32
var externalHost string
var name string
var reIndex bool

func init() {
	serveCmd.Flags().Int32VarP(&grpcPort, "port", "p", 9000, "Grpc Port")
	serveCmd.Flags().StringVarP(&externalHost, "external-host", "H", "localhost", "External host")
	serveCmd.Flags().StringVarP(&name, "name", "n", "localhost-daemon", "Daemon name")
	serveCmd.Flags().BoolVar(&reIndex, "re-index", false, "Launch a full re-indexation")

	rootCmd.AddCommand(serveCmd)
}
