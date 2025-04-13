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

package main

import (
	"context"

	"github.com/mrsimonemms/opensesame/apps/provider-local/cmd"
	"github.com/mrsimonemms/opensesame/apps/provider-local/pkg/config"
	"github.com/mrsimonemms/opensesame/apps/provider-local/pkg/database"
	"github.com/mrsimonemms/opensesame/packages/go-sdk/provider"
	"github.com/rs/zerolog/log"
)

const (
	name        = "provider-local"
	description = "Authenticate with local credentials against a username and password"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading database")
	}

	ctx := context.Background()

	if err := db.Connect(ctx); err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to DB")
	}

	authenticationCmd := cmd.New(cfg, db)

	provider.NewProviderServer(name, description, authenticationCmd)
}
