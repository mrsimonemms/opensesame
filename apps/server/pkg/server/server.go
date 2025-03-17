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
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/database"
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

func (s *Server) setupRouter() *Server {
	return s
}

func New(cfg *config.ServerConfig, db database.Driver) *Server {
	app := fiber.New(fiber.Config{
		AppName:               "cloud-native-auth",
		DisableStartupMessage: true,
	})

	s := &Server{
		app:    app,
		config: cfg,
		db:     db,
	}

	return s.setupRouter()
}
