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

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func App() *fiber.App {
	return fiber.New(fiber.Config{
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
}
