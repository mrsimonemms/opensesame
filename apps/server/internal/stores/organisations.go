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

package stores

import (
	"context"
	"fmt"

	"github.com/mrsimonemms/cloud-native-auth/apps/server/internal/database"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
)

type Organisations struct {
	cfg *config.ServerConfig
	db  database.Driver
}

func (o *Organisations) CheckSlugIsUnique(ctx context.Context, slug string, expectedOrgID *string) (bool, error) {
	org, err := o.db.GetOrgBySlug(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("error getting org by slug: %w", err)
	}

	if org == nil {
		// No org found
		return true, nil
	}

	if expectedOrgID != nil && org.ID == *expectedOrgID {
		// If found org matched the given org ID
		return true, nil
	}

	return false, nil
}

func NewOrganisationsStore(cfg *config.ServerConfig, db database.Driver) *Organisations {
	return &Organisations{
		cfg: cfg,
		db:  db,
	}
}
