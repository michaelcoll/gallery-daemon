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

package photo

import (
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/service"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/caller"
	"github.com/michaelcoll/gallery-daemon/internal/photo/infrastructure/infra_repository"
	"github.com/michaelcoll/gallery-daemon/internal/photo/presentation"
)

type Module struct {
	photoService    service.PhotoService
	registerService service.RegisterService
	c               presentation.PhotoController
}

func (m Module) GetPhotoService() *service.PhotoService {
	return &m.photoService
}

func (m Module) GetRegisterService() *service.RegisterService {
	return &m.registerService
}

func (m Module) GetController() *presentation.PhotoController {
	return &m.c
}

func NewForServe(param model.ServeParameters) Module {
	repository := infra_repository.New()
	registerCaller := caller.New(param)

	return Module{
		photoService:    service.New(repository),
		c:               presentation.New(repository, param),
		registerService: service.NewRegisterService(registerCaller, param),
	}
}
func NewForIndex() Module {
	repository := infra_repository.New()

	return Module{photoService: service.New(repository)}
}
