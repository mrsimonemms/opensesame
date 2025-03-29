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

	"github.com/mrsimonemms/opensesame/apps/server/pkg/config"
)

type ProviderAccount struct {
	ID             string            `json:"-"` // Represents the database ID
	Tokens         map[string]string `json:"tokens" example:"accessToken:this-is-an-access-token,refreshToken:this-is-the-refresh-token"`
	ProviderUserID string            `json:"providerUserId" example:"11223344"`
	EmailAddress   *string           `json:"emailAddress" example:"test@test.com"`
	Name           *string           `json:"name" example:"Test Testington"`
	Username       *string           `json:"username" example:"testtestington"`
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
