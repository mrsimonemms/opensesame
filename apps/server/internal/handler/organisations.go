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
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
)

// Create organisation godoc
// @Summary		Create organisation
// @Description Create new organisation
// @Tags		Organisations
// @Accept		json
// @Produce		json
// @Success		200	"@todo"
// @Failure		400 "Validation error"
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs [post]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationCreate(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"type": "create org",
	})
}

// List organisations godoc
// @Summary		List organisations
// @Description List all organisations for the user
// @Tags		Organisations
// @Accept		json
// @Produce		json
// @Success		200	{object}	models.Pagination[models.Organisation]
// @Failure		400 "Validation error"
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs [get]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationList(c *fiber.Ctx) error {
	return c.JSON(&models.Pagination[*models.Organisation]{
		Data:       []*models.Organisation{},
		Count:      0,
		Page:       1,
		TotalPages: 1,
		PerPage:    25,
		Total:      0,
	})
}

// Delete organisation godoc
// @Summary		Delete organisation
// @Description Delete organisation
// @Tags		Organisations
// @Accept		json
// @Produce		json
// @Param		orgID	path	string	true	"Organisation ID"
// @Success		204	"No response"
// @Failure		400 "Validation error"
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs/{orgID} [delete]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationDelete(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	return c.JSON(fiber.Map{
		"type":  "delete org",
		"orgID": orgID,
	})
}

// Get organisation godoc
// @Summary		Get organisation
// @Description Get organisation
// @Tags		Organisations
// @Accept		json
// @Produce		json
// @Param		orgID	path	string	true	"Organisation ID"
// @Success		200	{object}	models.Organisation
// @Failure		400 "Validation error"
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs/{orgID} [get]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationGet(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	return c.JSON(fiber.Map{
		"type":  "get org",
		"orgID": orgID,
	})
}

// Update organisation godoc
// @Summary		Update organisation
// @Description Update organisation
// @Tags		Organisations
// @Accept		json
// @Produce		json
// @Param		orgID	path	string	true	"Organisation ID"
// @Success		200	{object}	models.Organisation
// @Failure		400 "Validation error"
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs/{orgID} [patch]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationUpdate(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	return c.JSON(fiber.Map{
		"type":  "update org",
		"orgID": orgID,
	})
}
