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

	"github.com/mrsimonemms/opensesame/apps/server/pkg/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Organisation struct {
	ID          bson.ObjectID       `bson:"_id,omitempty"`
	Name        string              `bson:"name"`
	Slug        string              `bson:"slug"`
	Users       []*OrganisationUser `bson:"users"`
	CreatedDate time.Time           `bson:"createdDate"`
	UpdatedDate time.Time           `bson:"updatedDate"`
}

func PaginateUniqueUsers(o *models.Organisation) (filter []bson.M, err error) {
	uniqueUsers := map[string]string{}
	for _, u := range o.Users {
		uniqueUsers[u.UserID] = u.UserID
	}

	filter = []bson.M{}
	for u := range uniqueUsers {
		id, err := bson.ObjectIDFromHex(u)
		if err != nil {
			return nil, fmt.Errorf("error converting org's user id to bson object id: %w", err)
		}

		filter = append(filter, bson.M{"_id": id})
	}

	return filter, nil
}

func (o *Organisation) ToModel() *models.Organisation {
	m := &models.Organisation{
		Name:        o.Name,
		Slug:        o.Slug,
		Users:       make([]*models.OrganisationUser, 0),
		CreatedDate: o.CreatedDate,
		UpdatedDate: o.UpdatedDate,
	}

	for _, u := range o.Users {
		m.Users = append(m.Users, u.ToModel())
	}

	if !o.ID.IsZero() {
		m.ID = o.ID.Hex()
	}

	return m
}

func OrganisationToMongo(m *models.Organisation) (*Organisation, error) {
	o := &Organisation{
		Name:        m.Name,
		Slug:        m.Slug,
		Users:       []*OrganisationUser{},
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
	}

	for _, user := range m.Users {
		o.Users = append(o.Users, OrganisationUserToMongo(user))
	}

	if m.ID != "" {
		id, err := bson.ObjectIDFromHex(m.ID)
		if err != nil {
			return nil, fmt.Errorf("error converting organisation id to bson object id: %w", err)
		}

		o.ID = id
	}

	return o, nil
}
