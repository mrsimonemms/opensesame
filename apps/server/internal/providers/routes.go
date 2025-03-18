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
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
)

type controller struct {
	cfg *config.ServerConfig
}

func Router(route fiber.Router, cfg *config.ServerConfig) {
	p := controller{
		cfg: cfg,
	}

	route.Route("/providers", func(router fiber.Router) {
		router.Get("/", p.ListProviders)
	})
}

func (p *controller) ListProviders(c *fiber.Ctx) error {
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
