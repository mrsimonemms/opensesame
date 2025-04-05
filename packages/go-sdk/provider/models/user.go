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

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string    `json:"id" xml:"id" form:"name" example:"123456"` // Database
	Name         string    `json:"name" xml:"name" form:"name" example:"Test Testington" validate:"required"`
	Username     string    `json:"username" xml:"username" form:"username" example:"testtestington" validate:"required"`
	EmailAddress string    `json:"emailAddress" xml:"emailAddress" form:"emailAddress" example:"test@testington.com" validate:"required,email"`
	Password     string    `json:"password" xml:"password" form:"password" example:"this-is-some-password" validate:"required,min=8"`
	CreatedDate  time.Time `json:"createdDate" xml:"createdDate" form:"createdDate"`
	UpdatedDate  time.Time `json:"updatedDate" xml:"updatedDate" form:"updatedDate"`
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	u.Password = string(bytes)

	return nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}
