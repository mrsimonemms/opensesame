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
	"github.com/rs/zerolog/log"
)

func (h *handler) healthcheckProbe(c *fiber.Ctx) bool {
	if err := h.db.Check(c.Context()); err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		return false
	}

	log.Debug().Msg("Service healthy")
	return true
}
