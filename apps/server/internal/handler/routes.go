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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func (h *handler) Register(app *fiber.App) {
	app.
		Use(requestid.New()).
		Use(func(c *fiber.Ctx) error {
			l := log.With().
				Interface("requestid", c.Locals(requestid.ConfigDefault.ContextKey)).
				Str("method", c.Method()).
				Bytes("url", c.Request().URI().Path()). // Avoid logging any sensitive credentials
				Logger()

			c.Locals("logger", l)

			l.Debug().Msg("New route called")

			return c.Next()
		}).
		Use(recover.New()).
		Use(encryptcookie.New(encryptcookie.Config{
			Key: h.config.Server.Cookie.Key,
		}))

	// Health and observability checks
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe:  h.healthcheckProbe,
		ReadinessProbe: h.healthcheckProbe,
	}))
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Versioned endpoints
	v1 := app.Group("/v1")

	v1.Route("/providers", func(router fiber.Router) {
		router.Get("/", h.ProvidersList)
		router.Get(
			"/:providerID/login",
			h.IsRouteEnabled(authentication.Route_ROUTE_LOGIN_GET),
			h.VerifyUser(true),
			h.ProvidersLogin,
		)
		router.Post(
			"/:providerID/login",
			h.IsRouteEnabled(authentication.Route_ROUTE_LOGIN_POST),
			h.VerifyUser(true),
			h.ProvidersLogin,
		)
		router.Get("/:providerID/login/callback", h.IsRouteEnabled(authentication.Route_ROUTE_CALLBACK_GET), h.ProvidersLogin)
	})

	v1.Route("/user", func(router fiber.Router) {
		router.
			Use(h.VerifyUser()).
			Get("/", h.UserGet).
			Delete("/provider/:providerID", h.UserProviderDelete)
	})
}
