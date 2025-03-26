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

package handler

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/internal/providers"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (h *handler) IsRouteEnabled(route authentication.Route) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		providerID := c.Params("providerID")
		provider := providers.FindProvider(h.config.Providers, providerID)

		log := c.Locals("logger").(zerolog.Logger)

		l := log.With().Str("providerID", provider.ID).Str("route", route.String()).Logger()

		l.Debug().Msg("Validating route can be called")

		res, err := provider.Client.RouteEnabled(c.Context(), &authentication.RouteEnabledRequest{
			Route: route,
		})

		if err != nil || res.Enabled {
			if err != nil {
				l = l.With().Err(err).Logger()
			}
			l.Debug().Msg("Route enabled - continuing")
			return c.Next()
		}

		l.Debug().Msg("Route disabled")

		return fiber.ErrNotFound
	}
}

// Verify the user's permission to access the resource - errors with 403
func (h *handler) VerifyRBACPermissions(c *fiber.Ctx) error {
	return c.Next()
}

// Verifies the user's identity - errors with 401
func (h *handler) VerifyUser(isOptional ...bool) func(*fiber.Ctx) error {
	if len(isOptional) == 0 {
		isOptional = []bool{false}
	}

	return func(c *fiber.Ctx) error {
		var tokenLookup string
		token := c.Query(userAuthQueryString)
		if token != "" {
			// Token query string is set - authenticate through that
			tokenLookup = "query:" + userAuthQueryString
		}

		return jwtware.New(jwtware.Config{
			ContextKey:     jwtContextKey,
			ErrorHandler:   h.authErrorHandler(isOptional[0]),
			SuccessHandler: h.authSuccessHandler(isOptional[0]),
			SigningKey:     jwtware.SigningKey{Key: h.config.JWT.Key},
			TokenLookup:    tokenLookup,
		})(c)
	}
}

func (h *handler) authErrorHandler(isOptional bool) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		log.Debug().Err(err).Msg("Error validating user")

		return h.optionalErrorHandler(c, isOptional)
	}
}

func (h *handler) authSuccessHandler(isOptional bool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Locals(jwtContextKey).(*jwt.Token)

		now := time.Now()

		expiry, err := token.Claims.GetExpirationTime()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving expiry from JWT")
			return h.optionalErrorHandler(c, isOptional)
		}
		if expiry == nil || expiry.Before(now) {
			log.Debug().Msg("Token expiry invalid or expired")
			return h.optionalErrorHandler(c, isOptional)
		}

		notBefore, err := token.Claims.GetNotBefore()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving not before from JWT")
			return h.optionalErrorHandler(c, isOptional)
		}
		if notBefore == nil || notBefore.After(now) {
			log.Debug().Msg("Token not before invalid or expired")
			return h.optionalErrorHandler(c, isOptional)
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving issuer from JWT")
			return h.optionalErrorHandler(c, isOptional)
		}
		if issuer != h.config.JWT.Issuer {
			log.Debug().Msg("Token issuer invalid")
			return h.optionalErrorHandler(c, isOptional)
		}

		userID, err := token.Claims.GetSubject()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving user ID from JWT")
			return h.optionalErrorHandler(c, isOptional)
		}

		user, err := h.usersStore.GetUserByID(c.Context(), userID)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving user by ID")
			return h.optionalErrorHandler(c, isOptional)
		}
		if user == nil {
			log.Debug().Msg("No user found")
			return h.optionalErrorHandler(c, isOptional)
		}

		log.Debug().Msg("User found and saved to context")
		c.Locals(userContextKey, user)

		return c.Next()
	}
}

func (h *handler) optionalErrorHandler(c *fiber.Ctx, isOptional bool) error {
	if isOptional {
		return c.Next()
	}

	return fiber.ErrUnauthorized
}
