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

package providers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/config"
	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Convert the Fiber context to a Connect-like request object
func generateAuthRequest(c *fiber.Ctx) *authentication.AuthRequest {
	headers := map[string]*authentication.KeyRepeatedValue{}
	for key, value := range c.GetReqHeaders() {
		headers[strings.ToLower(key)] = &authentication.KeyRepeatedValue{
			Value: value,
		}
	}

	return &authentication.AuthRequest{
		Body:    string(c.Body()),
		Headers: headers,
		Method:  c.Method(),
		Query:   c.Queries(),
		Url:     c.OriginalURL(),
	}
}

// Handle authentication from the remove provider
func Authenticate(c *fiber.Ctx, provider config.Provider) (*authentication.User, error) {
	l := log.With().Str("providerID", provider.ID).Logger()

	l.Debug().Msg("Triggering call to gRPC provider")
	res, err := provider.Client.Auth(c.Context(), generateAuthRequest(c))
	if err != nil {
		statusCode, code, msg := ConvertGRPCErrorCodeToHTTP(err)

		l.Error().
			Err(err).
			Int("statusCode", statusCode).
			Uint32("grpcCode", uint32(code)).
			Str("errorMsg", msg).
			Msg("Error calling gRPC provider")

		return nil, fiber.NewError(statusCode, msg)
	}

	if res.Redirect != nil {
		l.Info().Int32("status", res.Redirect.Status).Msg("Auth redirecting")
		return nil, c.Redirect(res.Redirect.Url, int(res.Redirect.Status))
	}
	if res.Success != nil {
		l.Info().Msg("Auth successful")
		return res.Success.User, nil
	}

	l.Error().Msg("Empty AuthResponse received")
	return nil, fiber.ErrUnauthorized
}

func ConvertGRPCErrorCodeToHTTP(err error) (statusCode int, code codes.Code, msg string) {
	grpcError := status.Convert(err)
	code = grpcError.Code()
	msg = grpcError.Message()

	statusCode = fiber.StatusServiceUnavailable
	switch code {
	case codes.InvalidArgument, codes.FailedPrecondition:
		statusCode = fiber.StatusBadRequest
	case codes.NotFound:
		statusCode = fiber.StatusNotFound
	case codes.Unimplemented, codes.Internal:
		statusCode = fiber.StatusInternalServerError
	case codes.Unauthenticated:
		statusCode = fiber.StatusUnauthorized
	}

	return statusCode, code, msg
}

func FindProvider(providers []config.Provider, providerID string) *config.Provider {
	var provider *config.Provider
	for _, p := range providers {
		if p.ID == providerID {
			provider = &p
		}
	}
	return provider
}
