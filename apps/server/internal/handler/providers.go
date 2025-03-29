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
	"net/url"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mrsimonemms/opensesame/apps/server/internal/providers"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/models"
	"github.com/rs/zerolog"
)

type ProviderDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProviderLoginResponse struct {
	//nolint:lll // Allow long example
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4OTg1NzMxNTksImlhdCI6MTc0MjkwMzEyNiwiaXNzIjoiY2xvdWQtbmF0aXZlLWF1dGgiLCJuYmYiOjE3NDI5MDMxMjYsInN1YiI6IjUwN2YxZjc3YmNmODZjZDc5OTQzOTAxMSJ9.MfozqyuUj7pM8OX9JfYHyRu06JpcrioqBqYh5b8GlYI"`
	User  *models.User `json:"user"`
}

func (h *handler) ProvidersList(c *fiber.Ctx) error {
	providers := []ProviderDTO{}

	for _, i := range h.config.Providers {
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

// Login godoc
// @Summary		Login
// @Description Login to a provider
// @Tags		Providers
// @Accept		json
// @Produce		json
// @Success		200	{object}	ProviderLoginResponse
// @Response	302	"Redirect to provider login page"
// @Param		providerID	path	string	true	"Provider ID"	default(github)
// @Router		/v1/providers/{providerID}/login [get]
// @Router		/v1/providers/{providerID}/login [post]
// @Router		/v1/providers/{providerID}/login/callback [get]
// @Security	Token
func (h *handler) ProvidersLogin(c *fiber.Ctx) error {
	handleLoginInputCookies(c)

	log := c.Locals("logger").(zerolog.Logger)

	providerID := c.Params("providerID")
	provider := providers.FindProvider(h.config.Providers, providerID)
	if provider == nil {
		log.Debug().Str("providerID", providerID).Msg("Unknown provider ID")
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("Unknown provider ID: %s", providerID))
	}

	l := log.With().Str("providerID", providerID).Logger()

	l.Debug().Msg("Authenticating against provider")
	providerUser, err := providers.Authenticate(c, *provider)
	if err != nil {
		l.Error().Err(err).Msg("Error authenticating provider")
		return err
	}
	if providerUser == nil {
		// The webpage has successfully resolved - nothing to do
		return nil
	}

	l.Info().Msg("User authenticated by provider - saving to database")
	var existingUserID *string
	if userID := c.Cookies(existingUserCookieKey); userID != "" {
		l.Info().Str("existingUserID", userID).Msg("Existing user cookie found")
		existingUserID = &userID
	}

	l.Debug().Msg("Triggering user upsert")
	userModel, err := h.usersStore.CreateOrUpdateUserFromProvider(c.Context(), providerID, providerUser, existingUserID)
	if err != nil {
		l.Error().Err(err).Msg("Error creating user from provider")
		return fiber.NewError(fiber.StatusServiceUnavailable, "Error creating user from provider")
	}

	l.Debug().Msg("Generate the auth token")
	token, err := userModel.GenerateAuthToken(h.config)
	if err != nil {
		l.Error().Err(err).Msg("Error generating auth token")
		return fiber.NewError(fiber.StatusInternalServerError, "Error generating auth token")
	}

	if redirectURL := c.Cookies(callbackCookieKey); redirectURL != "" {
		l.Info().Msg("Redirecting to URL with token")

		u, err := url.Parse(redirectURL)
		if err != nil {
			l.Error().Err(err).Msg("Error parsing redirect URL")
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		q := u.Query()
		q.Add("token", token)
		u.RawQuery = q.Encode()

		return c.Redirect(u.String())
	}

	if err := userModel.DecryptTokens(h.config); err != nil {
		l.Error().Err(err).Msg("Error decrypting account tokens")
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	l.Info().Msg("Outputting the user object")
	return c.JSON(ProviderLoginResponse{
		Token: token,
		User:  userModel,
	})
}

func handleLoginInputCookies(c *fiber.Ctx) {
	log := c.Locals("logger").(zerolog.Logger)

	// Check if there is an existing user loaded
	if u := c.Locals(userContextKey); u != nil {
		existingUserID := (u.(*models.User)).ID
		log.Debug().Str("userID", existingUserID).Msg("Setting user ID to cookie")

		c.Cookie(&fiber.Cookie{
			Name:  existingUserCookieKey,
			Value: existingUserID,
		})
	} else {
		log.Debug().Msg("Clearing user ID cookie")
		// The c.ClearCookie function doesn't seem to work with encrypt cookie
		c.Cookie(&fiber.Cookie{
			Name:    existingUserCookieKey,
			Expires: time.Now().Add(-time.Hour * 24),
			Value:   "",
		})
	}

	if callbackURL := c.Query("callback", ""); callbackURL != "" {
		// Set a callback URL for after a success resolution
		log.Debug().Msg("Setting callback cookie")
		c.Cookie(&fiber.Cookie{
			Name:  callbackCookieKey,
			Value: callbackURL,
		})
	} else {
		log.Debug().Msg("Clearing callback cookie")
		// The c.ClearCookie function doesn't seem to work with encrypt cookie
		c.Cookie(&fiber.Cookie{
			Name:    callbackCookieKey,
			Expires: time.Now().Add(-time.Hour * 24),
			Value:   "",
		})
	}
}
