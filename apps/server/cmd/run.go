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

	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/database"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runOpts struct {
	ConfigFile string
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the service",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadFromFile(runOpts.ConfigFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Error loading config")
		}

		if err := cfg.Validate(); err != nil {
			log.Fatal().Err(err).Msg("Invalid config")
		}

		db, err := database.New(cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("Database connection error")
		}

		ctx := context.Background()

		log.Debug().Msg("Connecting to database")
		if err := db.Connect(ctx); err != nil {
			log.Fatal().Err(err).Msg("Error connecting to database")
		}

		defer db.Close(ctx)

		if err := server.New(cfg, db).Start(); err != nil {
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
