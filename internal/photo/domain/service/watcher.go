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

package service

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/consts"
	"log"

	"github.com/fsnotify/fsnotify"

	"github.com/michaelcoll/rfsnotify"
)

func (s *PhotoService) Watch(path string) {

	watcher, err := rfsnotify.NewBufferedWatcher(2000)
	if err != nil {
		log.Fatalf("Could not create the watcher : %v", err)
	}
	defer watcher.Close()

	err = watcher.AddRecursive(path, nil)
	if err != nil {
		log.Fatalf("Could not add the folder : %v", err)
	}

	fmt.Printf("%s Watching folder %s \n", color.GreenString("✓"), color.GreenString(path))

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if isCreateEvent(event) || isDeleteEvent(event) {
				s.handleEvent(event)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}

func isCreateEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Write
}

func isDeleteEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Rename ||
		event.Op == fsnotify.Remove
}

func (s *PhotoService) handleEvent(event fsnotify.Event) {
	if !hasExtension(event.Name, consts.SupportedExtensions) {
		return
	}

	s.r.Connect(false)
	defer s.r.Close()

	if isCreateEvent(event) {
		s.indexImage(context.Background(), event.Name)
	} else if isDeleteEvent(event) {
		s.deleteImage(context.Background(), event.Name)
	}
}
