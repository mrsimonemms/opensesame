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
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
	"github.com/rs/zerolog/log"
)

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

func Authenticate(c *fiber.Ctx, provider config.Provider) error {
	l := log.With().Str("providerID", provider.ID).Logger()

	l.Debug().Msg("Triggering call to gRPC provider")
	res, err := provider.Client.Auth(c.Context(), generateAuthRequest(c))
	if err != nil {
		l.Error().Err(err).Msg("Error calling gRPC provider")
		return fiber.ErrServiceUnavailable
	}

	fmt.Printf("%+v\n", res)

	return c.JSON(fiber.Map{
		"req": generateAuthRequest(c),
	})
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
