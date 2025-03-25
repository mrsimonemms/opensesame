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

package services

import (
	"context"
	"fmt"

	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/database"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
	"github.com/rs/zerolog/log"
)

type Users struct {
	cfg *config.ServerConfig
	db  database.Driver
}

func (s *Users) CreateOrUpdateUserFromProvider(
	ctx context.Context,
	providerID string,
	providerUser *authentication.User,
	existingUserID *string,
) (*models.User, error) {
	// Search for an existing user
	userModel, err := s.db.FindUserByProviderAndUserID(ctx, providerID, providerUser.ProviderId)
	if err != nil {
		return nil, fmt.Errorf("error getting user by provider and user id: %w", err)
	}

	if userModel == nil {
		log.Debug().Msg("No user found - creating")
		userModel = models.NewUser()

		// Add in default values
		if email := providerUser.EmailAddress; email != nil {
			userModel.EmailAddress = *email
		}
		if name := providerUser.Name; name != nil {
			userModel.Name = *name
		}
	}

	if existingUserID != nil {
		log.Info().Str("userID", *existingUserID).Msg("Linking provider to user")
		targetUser, err := s.db.GetUserByID(ctx, *existingUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting existing user by id: %w", err)
		}
		if targetUser == nil {
			return nil, fmt.Errorf("unknown user: %s", *existingUserID)
		}

		// Check if the tokens are used for a different account
		if userModel.ID == targetUser.ID {
			return nil, fmt.Errorf("provider registered with other user")
		}

		// Use the existing user from now on
		userModel = targetUser

		// Decode the tokens
		for _, a := range userModel.Accounts {
			if err := a.DecryptTokens(s.cfg); err != nil {
				log.Error().Err(err).Msg("Error decrypting account tokens")
				return nil, fmt.Errorf("error decrypting account tokens: %w", err)
			}
		}
	}

	userModel.AddProvider(providerID, providerUser)

	for _, a := range userModel.Accounts {
		if err := a.EncryptTokens(s.cfg); err != nil {
			log.Error().Err(err).Msg("Error encrypting account tokens")
			return nil, fmt.Errorf("error encrypting account tokens: %w", err)
		}
	}

	data, err := s.db.SaveUserRecord(ctx, userModel)
	if err != nil {
		return nil, fmt.Errorf("error saving user record: %w", err)
	}

	return data, nil
}

func (s *Users) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.db.GetUserByID(ctx, userID)
}

func NewUsersService(cfg *config.ServerConfig, db database.Driver) *Users {
	return &Users{
		cfg: cfg,
		db:  db,
	}
}
