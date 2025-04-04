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

	"github.com/mrsimonemms/opensesame/apps/provider-local/pkg/config"
	"github.com/mrsimonemms/opensesame/apps/provider-local/pkg/database/mongodb"
	"github.com/mrsimonemms/opensesame/packages/go-sdk/provider/models"
)

// Common database interface to allow multiple database types
// in the future.
type Driver interface {
	Check(ctx context.Context) error

	// Close the database connection and free up resources
	Close(ctx context.Context) error

	// Authorize the connection to the database
	Connect(ctx context.Context) error

	// Search for users by email address
	FindUserByEmailAddress(ctx context.Context, emailAddress string) (*models.User, error)

	// Search for users by username
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)

	// Save user
	Save(ctx context.Context, model *models.User) (*models.User, error)
}

func New(cfg *config.Config) (Driver, error) {
	var db Driver
	switch cfg.Database.Type {
	case "mongodb":
		db = mongodb.New(cfg.Database.MongoDB)
	default:
		return nil, fmt.Errorf("unknown database type")
	}

	return db, nil
}
