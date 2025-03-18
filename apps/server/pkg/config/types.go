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

type ServerConfig struct {
	Database  `json:"database" validate:"required"`
	Providers []Provider `json:"providers" validate:"required,min=1,dive"`
	Server    `json:"server" validate:"required"`
}

type DatabaseType string

const (
	DatabaseTypeMongoDB DatabaseType = "mongodb"
)

type Database struct {
	Type DatabaseType `json:"type" validate:"required,oneof=mongodb"`

	MongoDB `json:"mongodb" validate:"required_if=type mongodb"`
}

type MongoDB struct {
	ConnectionURI string `json:"connectionURI" validate:"required"`
	Database      string `json:"database" validate:"required"`
}

type Provider struct {
	ID      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required,hostname_port"`
}

type Server struct {
	Host string `json:"host" validate:"required,ip_addr"`
	Port int    `json:"port" validate:"required,number"`
}
