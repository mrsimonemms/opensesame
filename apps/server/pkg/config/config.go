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

package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"sigs.k8s.io/yaml"
)

func (s *ServerConfig) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("config failed validation: %w", err)
	}

	return nil
}

func LoadFromFile(configFile string) (*ServerConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	// Load the default values
	cfg := ServerConfig{
		Server: Server{
			Host: "0.0.0.0",
			Port: 3000,
		},
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal data")
	}

	return &cfg, nil
}
