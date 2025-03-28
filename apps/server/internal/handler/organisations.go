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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type OrgGetResponse struct {
	Org *models.Organisation `json:"org"`
}

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
// @Param		page	query	int	false	"Page number"	example(1)
// @Param		perPage	query	int	false	"Records per page"	example(25)
// @Success		200	{object}	models.Pagination[models.Organisation]
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs [get]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationList(c *fiber.Ctx) error {
	user := c.Locals(userContextKey).(*models.User)
	log := c.Locals("logger").(zerolog.Logger)

	page := max(c.QueryInt("page", 1), 1)
	perPage := min(max(c.QueryInt("perPage", 25), 1), 100)

	offset := perPage * (page - 1)

	log.Debug().
		Int("page", page).
		Int("perPage", perPage).
		Int("offset", offset).
		Str("userId", user.ID).
		Msg("Displaying organisations for user")

	orgs, err := h.db.ListOrganisations(c.Context(), offset, perPage, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting list of organisations")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(orgs)
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
// @Param		orgID	path	string	true	"Organisation ID" example(67e58132a5d5257f95a32518)
// @Success		200	{object}	OrgGetResponse
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs/{orgID} [get]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationGet(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	user := c.Locals(userContextKey).(*models.User)

	org, err := h.db.GetOrgByID(c.Context(), orgID, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting organisation")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if org == nil {
		return fiber.ErrNotFound
	}

	return c.JSON(OrgGetResponse{Org: org})
}

// List organisation's users godoc
// @Summary		List organisation's users
// @Description List all the users attached to an organisation.
// @Tags		Organisations
// @Accept		json
// @Produce		json
// @Param		orgID	path	string	true	"Organisation ID" example(67e58132a5d5257f95a32518)
// @Success		200	{object}	models.Pagination[OrganisationUser]
// @Failure		401 "Unauthorised error"
// @Router		/v1/orgs/{orgID}/users [get]
// @Security	Bearer
// @Security	Token
func (h *handler) OrganisationListUsers(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	user := c.Locals(userContextKey).(*models.User)
	log := c.Locals("logger").(zerolog.Logger)

	page := max(c.QueryInt("page", 1), 1)
	perPage := min(max(c.QueryInt("perPage", 25), 1), 100)

	offset := perPage * (page - 1)

	log.Debug().
		Int("page", page).
		Int("perPage", perPage).
		Int("offset", offset).
		Str("userId", user.ID).
		Str("orgId", orgID).
		Msg("Displaying users for organisation")

	org, err := h.db.ListOrganisationUsers(c.Context(), offset, perPage, orgID, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user for organisation")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if org == nil {
		return fiber.ErrNotFound
	}

	return c.JSON(org)
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
