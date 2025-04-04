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

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
	"github.com/mrsimonemms/opensesame/packages/go-sdk/provider/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Commands) UserCreate(
	ctx context.Context,
	request *authentication.UserCreateRequest,
) (*authentication.UserCreateResponse, error) {
	user := models.User{
		Name:         request.Name,
		EmailAddress: request.EmailAddress,
		Username:     request.Username,
		CreatedDate:  time.Now(),
		UpdatedDate:  time.Now(),
	}

	if err := user.SetPassword(request.Password); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.(validator.ValidationErrors).Error())
	}

	if user, err := c.db.FindUserByEmailAddress(ctx, user.EmailAddress); err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("error search for user by email: %w", err).Error())
	} else if user != nil {
		return nil, status.Error(codes.InvalidArgument, "email address already registered")
	}

	if user, err := c.db.FindUserByUsername(ctx, user.EmailAddress); err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("error search for user by username: %w", err).Error())
	} else if user != nil {
		return nil, status.Error(codes.InvalidArgument, "username already registered")
	}

	response, err := c.db.Save(ctx, &user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authentication.UserCreateResponse{
		Id:           response.ID,
		Name:         response.Name,
		Username:     response.Username,
		EmailAddress: response.EmailAddress,
	}, nil
}
