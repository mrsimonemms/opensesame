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

import "time"

type ProviderUser struct {
	ID             string            `json:"-"`
	Tokens         map[string]string `json:"tokens"` // This is highly sensitive
	ProviderID     string            `json:"providerId"`
	ProviderUserID string            `json:"providerUserId"`
	EmailAddress   *string           `json:"emailAddress"`
	Name           *string           `json:"name"`
	Username       *string           `json:"username"`
	CreatedDate    time.Time         `json:"createdDate"`
	UpdatedDate    time.Time         `json:"updatedDate"`
}

type User struct {
	ID           string          `json:"id"`
	EmailAddress string          `json:"emailAddress"`
	Name         string          `json:"name"`
	Accounts     []*ProviderUser `json:"accounts"`
	IsActive     bool            `json:"isActive"`
	CreatedDate  time.Time       `json:"createdDate"`
	UpdatedDate  time.Time       `json:"updatedDate"`
}
