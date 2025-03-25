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
	"fmt"
	"time"

	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
)

type ProviderAccount struct {
	ID             string            `json:"-"`      // Represents the database ID
	Tokens         map[string]string `json:"tokens"` // This is highly sensitive so will only be exported encrypted
	ProviderUserID string            `json:"providerUserId"`
	EmailAddress   *string           `json:"emailAddress"`
	Name           *string           `json:"name"`
	Username       *string           `json:"username"`
}

func (p *ProviderAccount) DecryptTokens(cfg *config.ServerConfig) error {
	for k, v := range p.Tokens {
		encrypted, err := decrypt(v, cfg.Encryption.Key)
		if err != nil {
			return fmt.Errorf("error decrypting provider token: %w", err)
		}

		p.Tokens[k] = string(encrypted)
	}

	return nil
}

func (p *ProviderAccount) EncryptTokens(cfg *config.ServerConfig) error {
	for k, v := range p.Tokens {
		encrypted, err := encrypt(v, cfg.Encryption.Key)
		if err != nil {
			return fmt.Errorf("error encrypting provider token: %w", err)
		}

		p.Tokens[k] = string(encrypted)
	}

	return nil
}

type User struct {
	ID           string                      `json:"id"` // Represents the database ID
	EmailAddress string                      `json:"emailAddress"`
	Name         string                      `json:"name"`
	Accounts     map[string]*ProviderAccount `json:"accounts"` // Key is the provider ID, eg github
	IsActive     bool                        `json:"isActive"`
	CreatedDate  time.Time                   `json:"createdDate"`
	UpdatedDate  time.Time                   `json:"updatedDate"`
}

func (u *User) AddProvider(providerID string, providerUser *authentication.User) {
	if u.Accounts[providerID] == nil {
		u.Accounts = map[string]*ProviderAccount{}
	}

	u.Accounts[providerID] = &ProviderAccount{
		Tokens:         providerUser.Tokens,
		ProviderUserID: providerUser.ProviderId,
		EmailAddress:   providerUser.EmailAddress,
		Name:           providerUser.Name,
		Username:       providerUser.Username,
	}

	u.UpdatedDate = time.Now()
}

func NewUser() *User {
	return &User{
		IsActive:    true, // Default to true
		CreatedDate: time.Now(),
	}
}
