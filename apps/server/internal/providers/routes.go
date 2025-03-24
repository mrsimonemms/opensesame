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
	"net/url"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/internal/services"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/auth"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/database"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
	"github.com/rs/zerolog"
)

type controller struct {
	cfg         *config.ServerConfig
	db          database.Driver
	userService *services.Users
}

const callbackCookieKey = "callback"

func Router(route fiber.Router, cfg *config.ServerConfig, db database.Driver) {
	p := controller{
		cfg:         cfg,
		db:          db,
		userService: services.NewUsersService(cfg, db),
	}

	route.Route("/providers", func(router fiber.Router) {
		router.Get("/", p.ListProviders)
		router.Get("/:providerID/login", p.IsRouteEnabled(authentication.Route_ROUTE_LOGIN_GET), p.LoginToProvider)
		router.Post("/:providerID/login", p.IsRouteEnabled(authentication.Route_ROUTE_LOGIN_POST), p.LoginToProvider)
		router.Get("/:providerID/login/callback", p.IsRouteEnabled(authentication.Route_ROUTE_CALLBACK_GET), p.LoginToProvider)
	})
}

func (p *controller) ListProviders(c *fiber.Ctx) error {
	providers := []ProviderDTO{}

	for _, i := range p.cfg.Providers {
		providers = append(providers, ProviderDTO{
			ID:   i.ID,
			Name: i.Name,
		})
	}

	// Put in alphabetical order
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Name < providers[j].Name
	})

	return c.JSON(providers)
}

func (p *controller) LoginToProvider(c *fiber.Ctx) error {
	log := c.Locals("logger").(zerolog.Logger)
	if callbackURL := c.Query("callback", ""); callbackURL != "" {
		// Set a callback URL for after a success resolution
		log.Debug().Msg("Setting callback cookie")
		c.Cookie(&fiber.Cookie{
			Name:  callbackCookieKey,
			Value: callbackURL,
		})
	} else {
		log.Debug().Msg("Clearing callback cookie")
		// The c.ClearCookie function doesn't seem to work with encrypt cookie
		c.Cookie(&fiber.Cookie{
			Name:    callbackCookieKey,
			Expires: time.Now().Add(-time.Hour * 24),
			Value:   "",
		})
	}

	providerID := c.Params("providerID")
	provider := FindProvider(p.cfg.Providers, providerID)
	if provider == nil {
		log.Debug().Str("providerID", providerID).Msg("Unknown provider ID")
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("Unknown provider ID: %s", providerID))
	}

	l := log.With().Str("providerID", providerID).Logger()

	l.Debug().Msg("Authenticating against provider")
	user, err := Authenticate(c, *provider)
	if err != nil {
		l.Error().Err(err).Msg("Error authenticating provider")
		return err
	}
	if user == nil {
		// The webpage has successfully resolved - nothing to do
		return nil
	}

	l.Info().Msg("User authenticated by provider - saving to database")

	l.Debug().Msg("Triggering user upsert")
	userModel, err := p.userService.CreateOrUpdateUserFromProvider(c.Context(), providerID, user)
	if err != nil {
		l.Error().Err(err).Msg("Error creating user from provider")
		return fiber.NewError(fiber.StatusServiceUnavailable, "Error creating user from provider")
	}

	l.Debug().Msg("Generate the auth token")
	token, err := auth.GenerateToken(userModel, p.cfg)
	if err != nil {
		l.Error().Err(err).Msg("Error generating auth token")
		return fiber.NewError(fiber.StatusInternalServerError, "Error generating auth token")
	}

	if redirectURL := c.Cookies(callbackCookieKey); redirectURL != "" {
		l.Info().Msg("Redirecting to URL with token")

		u, err := url.Parse(redirectURL)
		if err != nil {
			l.Error().Err(err).Msg("Error parsing redirect URL")
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		q := u.Query()
		q.Add("token", token)
		u.RawQuery = q.Encode()

		return c.Redirect(u.String())
	}

	l.Info().Msg("Outputting the user object")
	return c.JSON(fiber.Map{
		"token": token,
		"user":  userModel,
	})
}

func (p *controller) IsRouteEnabled(route authentication.Route) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		providerID := c.Params("providerID")
		provider := FindProvider(p.cfg.Providers, providerID)

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
