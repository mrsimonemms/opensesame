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
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
	"github.com/rs/zerolog"
)

func (h *handler) UserGet(c *fiber.Ctx) error {
	user := c.Locals(userContextKey).(*models.User)
	log := c.Locals("logger").(zerolog.Logger)

	for providerID, a := range user.Accounts {
		if err := a.DecryptTokens(h.config); err != nil {
			log.Error().Err(err).Str("providerID", providerID).Msg("Error decrypting provider token")
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Error decrypting provider token: %s", providerID))
		}
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func (h *handler) UserProviderDelete(c *fiber.Ctx) error {
	providerID := c.Params("providerID")
	user := c.Locals(userContextKey).(*models.User)
	log := c.Locals("logger").(zerolog.Logger)

	log = log.With().Str("providerID", providerID).Logger()

	_, err := h.usersStore.RemoveProviderFromUser(c.Context(), user.ID, providerID)
	if err != nil {
		log.Warn().Err(err).Msg("Error updating user")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
