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

package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
)

func GenerateToken(user *models.User, cfg *config.ServerConfig) (string, error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(cfg.ExpiresIn.Duration).Unix(),
			"iat": time.Now().Unix(),
			"iss": cfg.JWT.Issuer,
			"nbf": time.Now().Unix(),
			"sub": user.ID,
		},
	)

	s, err := t.SignedString(cfg.JWT.Key)
	if err != nil {
		return "", fmt.Errorf("error generating jwt signed string: %w", err)
	}

	return s, nil
}
