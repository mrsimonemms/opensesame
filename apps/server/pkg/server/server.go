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

package server

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/internal/providers"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/internal/users"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Server struct {
	app    *fiber.App
	config *config.ServerConfig
	db     database.Driver
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	log.Info().Str("address", addr).Msg("Starting server")

	return s.app.Listen(addr)
}

func (s *Server) healthcheckProbe(c *fiber.Ctx) bool {
	if err := s.db.Check(c.Context()); err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		return false
	}

	log.Debug().Msg("Service healthy")
	return true
}

func (s *Server) setupRouter() *Server {
	log.Debug().Msg("Creating routes")

	// Health and observability checks
	s.app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe:  s.healthcheckProbe,
		ReadinessProbe: s.healthcheckProbe,
	}))
	s.app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Versioned endpoints
	v1 := s.app.Group("/v1")
	providers.Router(v1, s.config, s.db)
	users.Router(v1, s.config, s.db)

	return s
}

func New(cfg *config.ServerConfig, db database.Driver) *Server {
	app := fiber.New(fiber.Config{
		AppName:               "cloud-native-auth",
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			if code >= 500 && code < 600 {
				// Log as developer error
				log.Error().Err(err).Msg("Error")
			} else {
				// Log as human error
				log.Debug().Err(err).Msg(e.Message)
			}

			// Render the error as JSON
			err = c.Status(code).JSON(e)
			if err != nil {
				log.Error().Err(err).Msg("Error rendering web output")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
			}

			return nil
		},
	})

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
			Key: cfg.Server.Cookie.Key,
		}))

	s := &Server{
		app:    app,
		config: cfg,
		db:     db,
	}

	return s.setupRouter()
}
