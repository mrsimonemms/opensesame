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

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Config struct {
	Database `envPrefix:"DATABASE_"`
	Logger   `envPrefix:"LOGGER_"`
}

type Database struct {
	Type string `env:"TYPE" envDefault:"mongodb" validate:"required,oneof=mongodb"`

	MongoDB `envPrefix:"MONGODB_" validate:"required_if=type mongodb"`
}

type Logger struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

type MongoDB struct {
	ConnectionURI string `env:"CONNECTION_URI" validate:"required"`
	Database      string `env:"DATABASE" validate:"required"`
	Collection    string `env:"COLLECTION" envDefault:"local_users" validate:"required"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config failed validation: %w", err)
	}

	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("error getting log level: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	return &cfg, nil
}
