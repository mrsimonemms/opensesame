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

package provider

import (
	"context"
	"fmt"
	"maps"
	"strings"

	grpcHelper "github.com/mrsimonemms/golang-helpers/grpc"
	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewProviderServer(name, description string, authenticationCmd authentication.AuthenticationServiceServer) {
	g := grpcHelper.New(name, description, []grpcHelper.ServerFactory{
		func(server *grpc.Server) {
			authentication.RegisterAuthenticationServiceServer(server, authenticationCmd)
		},
	})

	grpcHelper.NewGRPCCommand(g, "auth", grpcHelper.Listener[authentication.AuthResponse]{
		Flags: func(c *cobra.Command) {
			c.Flags().String("body", "{}", "Body input")
			c.Flags().StringToString("headers", map[string]string{}, "")
			c.Flags().String("method", "GET", "Method input")
			c.Flags().StringToString("query", map[string]string{}, "")
			c.Flags().String("url", "GET", "URL input")
		},
		Run: func(c *cobra.Command, s []string) (*authentication.AuthResponse, error) {
			// Receive the input values
			body, err := c.Flags().GetString("body")
			cobra.CheckErr(err)

			headersRaw, err := c.Flags().GetStringToString("headers")
			cobra.CheckErr(err)

			headers := map[string]*authentication.KeyRepeatedValue{}
			for key, value := range headersRaw {
				headers[key] = &authentication.KeyRepeatedValue{
					Value: strings.Split(value, ","),
				}
			}

			method, err := c.Flags().GetString("method")
			cobra.CheckErr(err)

			query, err := c.Flags().GetStringToString("query")
			cobra.CheckErr(err)

			url, err := c.Flags().GetString("url")
			cobra.CheckErr(err)

			return authenticationCmd.Auth(context.Background(), &authentication.AuthRequest{
				Body:    body,
				Headers: headers,
				Method:  method,
				Query:   query,
				Url:     url,
			})
		},
	})

	grpcHelper.NewGRPCCommand(g, "routeEnabled", grpcHelper.Listener[authentication.RouteEnabledResponse]{
		Flags: func(c *cobra.Command) {
			c.Flags().String("route", "", "Route name")
		},
		Run: func(c *cobra.Command, s []string) (*authentication.RouteEnabledResponse, error) {
			routeName, err := c.Flags().GetString("route")
			cobra.CheckErr(err)

			route, ok := authentication.Route_value[routeName]
			if !ok {
				keys := make([]string, 0)
				for k := range authentication.Route_value {
					keys = append(keys, k)
				}

				return nil, fmt.Errorf("unknown route, must be one of: %s", strings.Join(keys, ", "))
			}

			return authenticationCmd.RouteEnabled(context.Background(), &authentication.RouteEnabledRequest{
				Route: authentication.Route(route),
			})
		},
	})

	g.Execute()
}

func RoutesEnabled(
	request *authentication.RouteEnabledRequest,
	routeOverrides map[authentication.Route]bool,
) (*authentication.RouteEnabledResponse, error) {
	// Default everything to false
	routes := map[authentication.Route]bool{
		authentication.Route_ROUTE_CALLBACK_GET: false,
		authentication.Route_ROUTE_LOGIN_GET:    false,
		authentication.Route_ROUTE_LOGIN_POST:   false,
	}

	// Merge the overrides
	maps.Copy(routes, routeOverrides)

	enabled, ok := routes[request.Route]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "unknown route: %d", request.Route)
	}

	return &authentication.RouteEnabledResponse{
		Enabled: enabled,
	}, nil
}
