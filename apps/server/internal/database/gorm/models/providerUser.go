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

package models

import (
	"time"

	"github.com/google/uuid"
)

type ProviderUser struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Tokens         map[string]string
	ProviderUserID string
	EmailAddress   *string
	Name           *string
	Username       *string
	CreatedAt      time.Time `gorm:"column:createdDate"`
	UpdatedAt      time.Time `gorm:"column:updatedDate"`
}
