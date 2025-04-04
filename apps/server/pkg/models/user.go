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

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/config"
	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
)

type User struct {
	ID           string                      `json:"id" example:"507f1f77bcf86cd799439011"` // Represents the database ID
	EmailAddress string                      `json:"emailAddress" example:"test@test.com"`
	Name         string                      `json:"name" example:"Test Testington"`
	Accounts     map[string]*ProviderAccount `json:"accounts"` // Key is the provider ID, eg github
	IsActive     bool                        `json:"isActive" example:"true"`
	CreatedDate  time.Time                   `json:"createdDate" format:"date-time"`
	UpdatedDate  time.Time                   `json:"updatedDate" format:"date-time"`
}

func (u *User) AddProvider(providerID string, providerUser *authentication.User) {
	if u.Accounts == nil {
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

func (u *User) DecryptTokens(cfg *config.ServerConfig) error {
	for provider, accounts := range u.Accounts {
		if err := accounts.DecryptTokens(cfg); err != nil {
			return fmt.Errorf("error decrypting account tokens for %s: %w", provider, err)
		}
	}
	return nil
}

func (u *User) EncryptTokens(cfg *config.ServerConfig) error {
	for provider, accounts := range u.Accounts {
		if err := accounts.EncryptTokens(cfg); err != nil {
			return fmt.Errorf("error encrypting account tokens for %s: %w", provider, err)
		}
	}
	return nil
}

func (u *User) GenerateAuthToken(cfg *config.ServerConfig) (string, error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(cfg.ExpiresIn.Duration).Unix(),
			"iat": time.Now().Unix(),
			"iss": cfg.JWT.Issuer,
			"nbf": time.Now().Unix(),
			"sub": u.ID,
		},
	)

	s, err := t.SignedString(cfg.JWT.Key)
	if err != nil {
		return "", fmt.Errorf("error generating jwt signed string: %w", err)
	}

	return s, nil
}

func NewUser() *User {
	return &User{
		IsActive:    true, // Default to true
		CreatedDate: time.Now(),
	}
}
