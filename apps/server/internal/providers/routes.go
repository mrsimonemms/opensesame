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
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
	"github.com/rs/zerolog/log"
)

type providers struct {
	cfg *config.ServerConfig
}

func Router(route fiber.Router, cfg *config.ServerConfig) {
	p := providers{
		cfg: cfg,
	}

	route.Route("/providers", func(router fiber.Router) {
		router.Get("/", p.ListProviders)
		router.Get("/:providerID/login", p.ValidateRoute(authentication.Route_ROUTE_LOGIN_GET), p.LoginToProvider)
		router.Post("/:providerID/login", p.ValidateRoute(authentication.Route_ROUTE_LOGIN_POST), p.LoginToProvider)
		router.Get("/:providerID/login/callback", p.ValidateRoute(authentication.Route_ROUTE_CALLBACK_GET), p.LoginToProvider)
	})
}

func (p *providers) ListProviders(c *fiber.Ctx) error {
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

func (p *providers) ValidateRoute(routeID authentication.Route) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		providerID := c.Params("providerID")
		provider := FindProvider(p.cfg.Providers, providerID)

		l := log.With().Str("providerID", provider.ID).Int32("route", int32(routeID)).Logger()

		l.Debug().Msg("Validating route can be called")

		res, err := provider.Client.RouteEnabled(c.Context(), &authentication.RouteEnabledRequest{
			Route: authentication.Route_ROUTE_CALLBACK_GET,
		})
		fmt.Printf("%+v\n", routeID)
		fmt.Printf("%+v\n", authentication.Route_ROUTE_CALLBACK_GET)

		return fiber.ErrNotFound
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

func (p *providers) LoginToProvider(c *fiber.Ctx) error {
	providerID := c.Params("providerID")
	provider := FindProvider(p.cfg.Providers, providerID)
	if provider == nil {
		log.Debug().Str("providerID", providerID).Msg("Unknown provider ID")
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("Unknown provider ID: %s", providerID))
	}

	log.Debug().Str("providerID", providerID).Msg("Authenticating against provider")
	return Authenticate(c, *provider)
}
