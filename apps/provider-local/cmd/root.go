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

package cmd

import (
	"github.com/mrsimonemms/opensesame/apps/provider-local/pkg/config"
	"github.com/mrsimonemms/opensesame/apps/provider-local/pkg/database"
	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
)

type Commands struct {
	authentication.UnimplementedAuthenticationServiceServer

	cfg *config.Config
	db  database.Driver
}

func New(cfg *config.Config, db database.Driver) *Commands {
	return &Commands{
		cfg: cfg,
		db:  db,
	}
}
