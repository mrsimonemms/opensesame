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

//nolint:misspell
package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Server struct {
	app *fiber.App
}

func (s *Server) Start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	log.Debug().Str("address", addr).Msg("Starting server")

	return s.app.Listen(addr)
}

func New() *Server {
	app := fiber.New()

	return &Server{
		app: app,
	}
}
