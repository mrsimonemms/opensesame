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

package auth

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/internal/services"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/database"
	"github.com/rs/zerolog/log"
)

const (
	jwtContextKey       = "jwtoken"
	userAuthQueryString = "token"
	UserContextKey      = "user"
)

// func OptionallyVerifyUser(cfg *config.ServerConfig, db database.Driver) func(*fiber.Ctx) error {
// 	return func(c *fiber.Ctx) error {
// 		log := c.Locals("logger").(zerolog.Logger)

// 		log.Debug().Msg("Optionally checking for authenticated user")

// 		err := VerifyUser(cfg, db)(c)
// 		if err != nil {
// 			log.Debug().Err(err).Msg("User not verified - continuing as anonymous")
// 		} else {
// 			log.Debug().Msg("User verified - continuing as authenticated")
// 		}

// 		fmt.Println("triggering next")
// 		return c.Next()
// 	}
// }

// Verify the user's permission to access the resource - errors with 403
func VerifyRBACPermissions(cfg *config.ServerConfig, db database.Driver) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// Verifies the user's identity - errors with 401
func VerifyUser(cfg *config.ServerConfig, db database.Driver, isOptional ...bool) func(*fiber.Ctx) error {
	if len(isOptional) == 0 {
		isOptional = []bool{false}
	}

	usersService := services.NewUsersService(cfg, db)

	return func(c *fiber.Ctx) error {
		var tokenLookup string
		token := c.Query(userAuthQueryString)
		if token != "" {
			// Token query string is set - authenticate through that
			tokenLookup = "query:" + userAuthQueryString
		}

		return jwtware.New(jwtware.Config{
			ContextKey:     jwtContextKey,
			ErrorHandler:   authErrorHandler(isOptional[0]),
			SuccessHandler: authSuccessHandler(cfg, usersService, isOptional[0]),
			SigningKey:     jwtware.SigningKey{Key: cfg.JWT.Key},
			TokenLookup:    tokenLookup,
		})(c)
	}
}

func optionalErrorHandler(c *fiber.Ctx, isOptional bool) error {
	if isOptional {
		return c.Next()
	}

	return fiber.ErrUnauthorized
}

func authErrorHandler(isOptional bool) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		log.Debug().Err(err).Msg("Error validating user")

		return optionalErrorHandler(c, isOptional)
	}
}

func authSuccessHandler(cfg *config.ServerConfig, usersService *services.Users, isOptional bool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Locals(jwtContextKey).(*jwt.Token)

		now := time.Now()

		expiry, err := token.Claims.GetExpirationTime()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving expiry from JWT")
			return optionalErrorHandler(c, isOptional)
		}
		if expiry == nil || expiry.Before(now) {
			log.Debug().Msg("Token expiry invalid or expired")
			return optionalErrorHandler(c, isOptional)
		}

		notBefore, err := token.Claims.GetNotBefore()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving not before from JWT")
			return optionalErrorHandler(c, isOptional)
		}
		if notBefore == nil || notBefore.After(now) {
			log.Debug().Msg("Token not before invalid or expired")
			return optionalErrorHandler(c, isOptional)
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving issuer from JWT")
			return optionalErrorHandler(c, isOptional)
		}
		if issuer != cfg.JWT.Issuer {
			log.Debug().Msg("Token issuer invalid")
			return optionalErrorHandler(c, isOptional)
		}

		userID, err := token.Claims.GetSubject()
		if err != nil {
			log.Debug().Err(err).Msg("Error retrieving user ID from JWT")
			return optionalErrorHandler(c, isOptional)
		}

		user, err := usersService.GetUserByID(c.Context(), userID)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving user by ID")
			return optionalErrorHandler(c, isOptional)
		}
		if user == nil {
			log.Debug().Msg("No user found")
			return optionalErrorHandler(c, isOptional)
		}

		log.Debug().Msg("User found and saved to context")
		c.Locals(UserContextKey, user)

		return c.Next()
	}
}
