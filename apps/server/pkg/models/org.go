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
)

type OrganisationUser struct {
	ID           string    `json:"-"`
	UserID       string    `json:"userId" form:"userId" example:"507f1f77bcf86cd799439011" validate:"required"`
	Name         string    `json:"name,omitempty" example:"Test Testington"`
	EmailAddress string    `json:"emailAddress,omitempty" example:"test@testington.com"`
	Role         string    `json:"role" form:"role" example:"ORG_MAINTAINER" validate:"required"`
	CreatedDate  time.Time `json:"createdDate" format:"date-time"`
	UpdatedDate  time.Time `json:"updatedDate" format:"date-time"`
}

type Organisation struct {
	ID          string              `json:"id" example:"67e58132a5d5257f95a32518"` // Represents the database ID
	Name        string              `json:"name" form:"name" example:"Org Name" validate:"required"`
	Slug        string              `json:"slug" form:"slug" example:"orgname" validate:"required"`
	Users       []*OrganisationUser `json:"users" form:"users" validate:"required"`
	CreatedDate time.Time           `json:"createdDate" format:"date-time"`
	UpdatedDate time.Time           `json:"updatedDate" format:"date-time"`
}

type OrgDTO struct {
	Name  string        `json:"name" form:"name" example:"Org Name" validate:"required"`
	Slug  string        `json:"slug" form:"slug" example:"orgname" validate:"required"`
	Users []*OrgUserDTO `json:"users" form:"users" validate:"required,min=1"`
}

func (o *OrgDTO) ToModel() *Organisation {
	m := &Organisation{
		Name:  o.Name,
		Slug:  o.Slug,
		Users: []*OrganisationUser{},
	}

	for _, u := range o.Users {
		m.Users = append(m.Users, &OrganisationUser{
			UserID: u.UserID,
			Role:   u.Role,
		})
	}

	return m
}

type OrgUserDTO struct {
	UserID string `json:"userId" example:"507f1f77bcf86cd799439011"`
	Role   string `json:"role" example:"ORG_MAINTAINER"`
}
