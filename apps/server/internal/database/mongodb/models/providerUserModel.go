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

import "github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"

type ProviderUser struct {
	Tokens         map[string]string `bson:"tokens"`
	ProviderUserID string            `bson:"providerUserId"`
	EmailAddress   *string           `bson:"emailAddress"`
	Name           *string           `bson:"name"`
	Username       *string           `bson:"username"`
}

func (p *ProviderUser) ToModel() *models.ProviderAccount {
	return &models.ProviderAccount{
		Tokens:         p.Tokens,
		ProviderUserID: p.ProviderUserID,
		EmailAddress:   p.EmailAddress,
		Name:           p.Name,
		Username:       p.Username,
	}
}

func ProviderUserToMongo(p *models.ProviderAccount) *ProviderUser {
	return &ProviderUser{
		Tokens:         p.Tokens,
		ProviderUserID: p.ProviderUserID,
		EmailAddress:   p.EmailAddress,
		Name:           p.Name,
		Username:       p.Username,
	}
}
