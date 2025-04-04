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

	"github.com/mrsimonemms/opensesame/packages/go-sdk/provider/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	EmailAddress string        `bson:"emailAddress"`
	Name         string        `bson:"name"`
	Username     string        `bson:"username"`
	Password     string        `bson:"password"`
	CreatedDate  time.Time     `bson:"createdDate"`
	UpdatedDate  time.Time     `bson:"updatedDate"`
}

func (u *User) ToModel() *models.User {
	m := &models.User{
		EmailAddress: u.EmailAddress,
		Name:         u.Name,
		Username:     u.Username,
		Password:     u.Password,
		CreatedDate:  u.CreatedDate,
		UpdatedDate:  u.UpdatedDate,
	}

	if !u.ID.IsZero() {
		m.ID = u.ID.Hex()
	}

	return m
}

func UserToModel(m *models.User) (*User, error) {
	u := &User{
		EmailAddress: m.EmailAddress,
		Name:         m.Name,
		Username:     m.Username,
		Password:     m.Password,
		CreatedDate:  m.CreatedDate,
		UpdatedDate:  m.UpdatedDate,
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
