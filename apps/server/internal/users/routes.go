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

package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/auth"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/database"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
	"github.com/rs/zerolog"
)

type controller struct {
	cfg *config.ServerConfig
	db  database.Driver
}

func Router(route fiber.Router, cfg *config.ServerConfig, db database.Driver) {
	p := controller{
		cfg: cfg,
		db:  db,
	}

	route.Route("/user", func(router fiber.Router) {
		router.Get("/", auth.VerifyUser(cfg, db), p.GetUser)
	})
}

func (p *controller) GetUser(c *fiber.Ctx) error {
	user := c.Locals(auth.UserContextKey).(*models.User)
	log := c.Locals("logger").(zerolog.Logger)

	for _, a := range user.Accounts {
		if err := a.DecryptTokens(p.cfg); err != nil {
			log.Error().Err(err).Str("providerID", a.ProviderID).Msg("Error decrypting provider token")
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Error decrypting provider token: %s", a.ProviderID))
		}
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}
