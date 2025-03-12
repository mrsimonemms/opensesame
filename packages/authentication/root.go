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

package authentication

import (
	"context"

	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
	grpcHelper "github.com/mrsimonemms/golang-helpers/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func New(name, description string, authCmd authentication.AuthenticationServiceServer) *grpcHelper.Server {
	g := grpcHelper.New(name, description, []grpcHelper.ServerFactory{
		func(server *grpc.Server) {
			authentication.RegisterAuthenticationServiceServer(server, authCmd)
		},
	})

	// Define the auth command
	grpcHelper.NewGRPCCommand(g, "auth", grpcHelper.Listener[authentication.AuthResponse]{
		Run: func(c *cobra.Command, s []string) (*authentication.AuthResponse, error) {
			return authCmd.Auth(context.Background(), &authentication.AuthRequest{})
		},
	})

	return g
}
