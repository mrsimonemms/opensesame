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
	"time"

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
) (*models.User, error) {
	userModel, err := s.db.FindUserByProviderAndUserID(ctx, providerID, providerUser.ProviderId)
	if err != nil {
		return nil, fmt.Errorf("error getting user by provider and user id: %w", err)
	}

	now := time.Now()

	if userModel == nil {
		log.Debug().Msg("No user found - creating")
		userModel = &models.User{
			IsActive:    true, // Default to true
			CreatedDate: now,
		}

		// Add email and name from the provider user - it's not guaranteed to be present
		if providerUser.EmailAddress != nil {
			userModel.EmailAddress = *providerUser.EmailAddress
		}
		if providerUser.Name != nil {
			userModel.Name = *providerUser.Name
		}
	}

	// Search for provider account in slice
	var accountID *int
	for k, account := range userModel.Accounts {
		if account.ProviderID == providerID && account.ProviderUserID == providerUser.ProviderId {
			accountID = &k
		}
	}

	accountRecord := models.ProviderUser{
		Tokens:         providerUser.Tokens,
		ProviderID:     providerID,
		ProviderUserID: providerUser.ProviderId,
		EmailAddress:   providerUser.EmailAddress,
		Name:           providerUser.Name,
		Username:       providerUser.Username,
		UpdatedDate:    now,
	}

	if accountID == nil {
		log.Debug().Msg("Adding new account record")
		accountRecord.CreatedDate = now
		userModel.Accounts = append(userModel.Accounts, &accountRecord)
	} else {
		log.Debug().Int("accountID", *accountID).Msg("Updating account record")

		// Set the extant records
		accountRecord.ID = userModel.Accounts[*accountID].ID
		accountRecord.CreatedDate = userModel.Accounts[*accountID].CreatedDate

		userModel.Accounts[*accountID] = &accountRecord
	}
	userModel.UpdatedDate = now

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
