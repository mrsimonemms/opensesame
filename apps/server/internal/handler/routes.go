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

//go:generate swag init --output ../../docs -g routes.go --parseDependency --parseInternal

package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	_ "github.com/mrsimonemms/opensesame/apps/server/docs"
)

// @title Open Sesame
// @version 1.0
// @description Authentication and authorisation for cloud-native apps
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
// @contact.name Open Sesame
// @contact.url https://opensesame.cloud
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
//
// @securityDefinitions.apikey Token
// @in query
// @name token
// @description Type JWT token.
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

	app.Get("api/*", swagger.HandlerDefault)

	// Health and observability checks
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe:  h.healthcheckProbe,
		ReadinessProbe: h.healthcheckProbe,
	}))
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Versioned endpoints
	v1 := app.Group("/v1")

	v1.Route("/orgs", func(router fiber.Router) {
		router.
			Use(h.VerifyUser()).
			Get("/", h.OrganisationList).
			Post("/", h.OrganisationCreate)

		router.Route("/:orgID", func(r fiber.Router) {
			r.
				Get("/", h.OrganisationGet).
				Delete("/", h.OrganisationDelete).
				Patch("/", h.OrganisationUpdate).
				Get("/users", h.OrganisationListUsers)
		})
	})

	v1.Route("/providers", func(router fiber.Router) {
		router.Get("/", h.ProvidersList)
		router.Get(
			"/:providerID/login",
			h.IsRouteEnabled(authentication.Route_ROUTE_LOGIN_GET),
			h.VerifyUser(true),
			h.ProvidersLogin,
		)
		router.Post(
			":providerID",
			h.IsRouteEnabled(authentication.Route_ROUTE_USER_CREATE),
			h.ProvidersCreateUser,
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
