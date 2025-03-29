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

package database

import (
	"context"
	"fmt"

	"github.com/mrsimonemms/opensesame/apps/server/internal/database/mongodb"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/config"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/models"
)

// Common database interface to allow multiple database types
// in the future.
type Driver interface {
	Check(ctx context.Context) error

	// Close the database connection and free up resources
	Close(ctx context.Context) error

	// Authorize the connection to the database
	Connect(ctx context.Context) error

	// Delete organisation
	DeleteOrganisation(ctx context.Context, orgID, userID string) error

	// Find the user by the provider and provider user ID
	FindUserByProviderAndUserID(ctx context.Context, providerID, providerUserID string) (user *models.User, err error)

	// Get the organisation by ID
	GetOrgByID(ctx context.Context, orgID, userID string) (org *models.Organisation, err error)

	// Get the organisation by Slug
	GetOrgBySlug(ctx context.Context, slug string) (org *models.Organisation, err error)

	// Get the user by ID
	GetUserByID(ctx context.Context, userID string) (user *models.User, err error)

	// List organisations available to a user
	ListOrganisations(ctx context.Context, offset, limit int, userID string) (orgs *models.Pagination[*models.Organisation], err error)

	// List users attached to an organisation
	ListOrganisationUsers(
		ctx context.Context,
		offset,
		limit int,
		orgID,
		userID string,
	) (users *models.Pagination[*models.OrganisationUser], err error)

	// Save the org record to the database
	SaveOrganisationRecord(ctx context.Context, model *models.Organisation) (user *models.Organisation, err error)

	// Save the user record to the database
	SaveUserRecord(ctx context.Context, model *models.User) (user *models.User, err error)

	// Updates all users - used when rotating keys
	UpdateAllUsers(ctx context.Context, update func(existing []*models.User) (updated []*models.User, err error)) (count int64, err error)
}

func New(cfg *config.ServerConfig) (Driver, error) {
	var db Driver
	switch cfg.Database.Type {
	case config.DatabaseTypeMongoDB:
		db = mongodb.New(cfg.Database.MongoDB)
	default:
		return nil, fmt.Errorf("unknown database type")
	}

	return db, nil
}
