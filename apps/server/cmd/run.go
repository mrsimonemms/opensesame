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
	"context"
	"fmt"

	"github.com/mrsimonemms/opensesame/apps/server/internal/database"
	"github.com/mrsimonemms/opensesame/apps/server/internal/handler"
	"github.com/mrsimonemms/opensesame/apps/server/internal/server"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runOpts struct {
	ConfigFile string
}

func loadConfig(configFile string) *config.ServerConfig {
	cfg, err := config.LoadFromFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid config")
	}

	if err := cfg.ConnectProviders(); err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to providers")
	}

	return cfg
}

func connectToDatabase(ctx context.Context, cfg *config.ServerConfig) database.Driver {
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Database connection error")
	}

	log.Debug().Msg("Connecting to database")
	if err := db.Connect(ctx); err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	return db
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the service",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cfg := loadConfig(runOpts.ConfigFile)
		db := connectToDatabase(ctx, cfg)

		defer db.Close(ctx)

		app := server.App()
		h := handler.New(cfg, db)
		h.Register(app)

		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

		log.Info().Str("address", addr).Msg("Starting server")

		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("Unable to start server")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(
		&runOpts.ConfigFile,
		"config",
		"c",
		bindEnv[string]("config", "config.yaml"),
		"Location to the config file",
	)
}
