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

	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
	sdk "github.com/mrsimonemms/opensesame/packages/go-sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/yaml"
)

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Commands) Auth(ctx context.Context, request *authentication.AuthRequest) (*authentication.AuthResponse, error) {
	// Convert the body JSON to the AuthBody request
	var body AuthBody
	if err := yaml.Unmarshal([]byte(request.Body), &body); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error parsing request body to yaml: %s", err.Error())
	}

	// If success is empty, authentication fails
	var success *authentication.Success

	if body.Username == "user" && body.Password == "valid" {
		success = &authentication.Success{
			User: &authentication.User{
				ProviderId:   "12345",
				Name:         sdk.Ptr("Test Testington"),
				Username:     &body.Username,
				EmailAddress: sdk.Ptr("test@testington.com"),
			},
		}
	}

	return &authentication.AuthResponse{
		Success: success,
	}, nil
}
