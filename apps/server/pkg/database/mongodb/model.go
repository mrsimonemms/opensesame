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

package mongodb

import (
	"fmt"
	"time"

	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Provider user-mapped model

type providerUser struct {
	Tokens         map[string]string `bson:"tokens"`
	ProviderID     string            `bson:"providerId"`
	ProviderUserID string            `bson:"providerUserId"`
	EmailAddress   *string           `bson:"emailAddress"`
	Name           *string           `bson:"name"`
	Username       *string           `bson:"username"`
	CreatedDate    time.Time         `bson:"createdDate"`
	UpdatedDate    time.Time         `bson:"updatedDate"`
}

func (p *providerUser) toModel() *models.ProviderUser {
	return &models.ProviderUser{
		Tokens:         p.Tokens,
		ProviderID:     p.ProviderID,
		ProviderUserID: p.ProviderUserID,
		EmailAddress:   p.EmailAddress,
		Name:           p.Name,
		Username:       p.Username,
		CreatedDate:    p.CreatedDate,
		UpdatedDate:    p.UpdatedDate,
	}
}

func providerUserToMongo(p *models.ProviderUser) *providerUser {
	return &providerUser{
		Tokens:         p.Tokens,
		ProviderID:     p.ProviderID,
		ProviderUserID: p.ProviderUserID,
		EmailAddress:   p.EmailAddress,
		Name:           p.Name,
		Username:       p.Username,
		CreatedDate:    p.CreatedDate,
		UpdatedDate:    p.UpdatedDate,
	}
}

// User-mapped model

type user struct {
	ID           bson.ObjectID   `bson:"_id,omitempty"`
	EmailAddress string          `bson:"emailAddress"`
	Name         string          `bson:"name"`
	Accounts     []*providerUser `bson:"accounts"`
	IsActive     bool            `bson:"isActive"`
	CreatedDate  time.Time       `bson:"createdDate"`
	UpdatedDate  time.Time       `bson:"updatedDate"`
}

func (u *user) toModel() *models.User {
	m := &models.User{
		EmailAddress: u.EmailAddress,
		Name:         u.Name,
		Accounts:     []*models.ProviderUser{},
		IsActive:     u.IsActive,
		CreatedDate:  u.CreatedDate,
		UpdatedDate:  u.UpdatedDate,
	}

	for _, i := range u.Accounts {
		m.Accounts = append(m.Accounts, i.toModel())
	}

	if !u.ID.IsZero() {
		m.ID = u.ID.Hex()
	}

	return m
}

func userToMongo(m *models.User) (*user, error) {
	u := &user{
		EmailAddress: m.EmailAddress,
		Name:         m.Name,
		Accounts:     []*providerUser{},
		IsActive:     m.IsActive,
		CreatedDate:  m.CreatedDate,
		UpdatedDate:  m.UpdatedDate,
	}

	for _, i := range m.Accounts {
		u.Accounts = append(u.Accounts, providerUserToMongo(i))
	}

	if m.ID != "" {
		id, err := bson.ObjectIDFromHex(m.ID)
		if err != nil {
			return nil, fmt.Errorf("error converting user id to bson object id: %w", err)
		}

		u.ID = id
	}

	return u, nil
}
