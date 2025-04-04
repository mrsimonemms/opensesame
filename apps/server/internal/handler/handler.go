/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/mrsimonemms/opensesame/apps/server/internal/database"
	"github.com/mrsimonemms/opensesame/apps/server/internal/stores"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/config"
)

type handler struct {
	config    *config.ServerConfig
	db        database.Driver
	validator *validator.Validate

	orgsStore  *stores.Organisations
	usersStore *stores.Users
}

func New(config *config.ServerConfig, db database.Driver) *handler {
	return &handler{
		config:    config,
		db:        db,
		validator: validator.New(validator.WithRequiredStructEnabled()),

		orgsStore:  stores.NewOrganisationsStore(config, db),
		usersStore: stores.NewUsersStore(config, db),
	}
}
