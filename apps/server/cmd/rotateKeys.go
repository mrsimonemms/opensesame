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

	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rotateKeysOpts struct {
	ConfigFile string
	NewKey     string
}

// rotateKeysCmd represents the rotateKeys command
var rotateKeysCmd = &cobra.Command{
	Use:     "rotateKeys",
	Aliases: []string{"rotate", "rotate-keys"},
	Short:   "Rotate encryption keys and update user account tokens",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cfg := loadConfig(rotateKeysOpts.ConfigFile)

		// Clone the config
		newCfg := *cfg
		newCfg.Encryption.Key = rotateKeysOpts.NewKey
		if err := newCfg.Validate(); err != nil {
			log.Fatal().Err(err).Msg("New config is invalid")
		}

		db := connectToDatabase(ctx, cfg)
		defer db.Close(ctx)

		updateCount, err := db.UpdateAllUsers(ctx, func(existing []*models.User) ([]*models.User, error) {
			for _, record := range existing {
				for _, a := range record.Accounts {
					if err := a.DecryptTokens(cfg); err != nil {
						return nil, fmt.Errorf("error decrypting account tokens: %w", err)
					}

					if err := a.EncryptTokens(&newCfg); err != nil {
						return nil, fmt.Errorf("error encrypting account tokens: %w", err)
					}
				}
			}

			return existing, nil
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error updating all user records")
		}

		log.Info().Int64("records updated", updateCount).Msg("User records updated with new key")
	},
}

func init() {
	rootCmd.AddCommand(rotateKeysCmd)

	rotateKeysCmd.Flags().StringVarP(
		&rotateKeysOpts.ConfigFile,
		"config",
		"c",
		bindEnv[string]("config", "config.yaml"),
		"Location to the config file",
	)
	rotateKeysCmd.Flags().StringVarP(&rotateKeysOpts.NewKey, "new-key", "k", bindEnv[string]("new-key", ""), "New encryption key")
}
