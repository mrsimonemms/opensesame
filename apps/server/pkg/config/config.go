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
	"bytes"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mrsimonemms/opensesame/packages/authentication/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sigs.k8s.io/yaml"
)

func (s *ServerConfig) ConnectProviders() error {
	for k, p := range s.Providers {
		log.Debug().Str("address", p.Address).Msg("Connecting to gRPC service")
		conn, err := grpc.NewClient(p.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("error connecting to grpc client: %w", err)
		}

		// Store the client for later use
		s.Providers[k].Client = authentication.NewAuthenticationServiceClient(conn)
	}
	return nil
}

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

	// Get the desired environment variables
	envvars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		if len(pair) == 2 && strings.HasPrefix(pair[0], EnvVarPrefix) {
			envvars[pair[0]] = pair[1]
		}
	}

	if len(envvars) > 0 {
		// Load envvars via Go templates
		log.Debug().Msg("Parsing config to include envvars")
		tpl, err := template.New("config").Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("error parsing config as template: %w", err)
		}

		// Execute the template
		var cfgParsed bytes.Buffer
		if err := tpl.Execute(&cfgParsed, envvars); err != nil {
			return nil, fmt.Errorf("error executing envvar template: %w", err)
		}

		data = cfgParsed.Bytes()
	}

	// Load the default values
	cfg := ServerConfig{
		JWT: JWT{
			ExpiresIn: Duration{
				Duration: time.Hour * 24 * 30, // 30 days,
			},
			Issuer: "opensesame.cloud",
		},
		Server: Server{
			Host: "0.0.0.0",
			Port: 3000,
		},
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal data: %w", err)
	}

	cfg.Providers = slices.DeleteFunc(cfg.Providers, func(p Provider) bool {
		return p.Disabled
	})

	return &cfg, nil
}

// ParseDuration parses a duration string.
// examples: "10d", "-1.5w" or "3Y4M5d".
// Add time units are "d"="D", "w"="W", "M", "y"="Y".
// @link https://gist.github.com/xhit/79c9e137e1cfe332076cdda9f5e24699
func ParseDuration(s string) (time.Duration, error) {
	neg := false
	if s != "" && s[0] == '-' {
		neg = true
		s = s[1:]
	}

	re := regexp.MustCompile(`(\d*\.\d+|\d+)\D*`)
	unitMap := map[string]time.Duration{
		"d": 24,
		"D": 24,
		"w": 7 * 24,
		"W": 7 * 24,
		"M": 30 * 24,
		"y": 365 * 24,
		"Y": 365 * 24,
	}

	strs := re.FindAllString(s, -1)
	var sumDur time.Duration
	for _, str := range strs {
		var _hours time.Duration = 1
		for unit, hours := range unitMap {
			if strings.Contains(str, unit) {
				str = strings.ReplaceAll(str, unit, "h")
				_hours = hours
				break
			}
		}

		dur, err := time.ParseDuration(str)
		if err != nil {
			return 0, err
		}

		sumDur += dur * _hours
	}

	if neg {
		sumDur = -sumDur
	}
	return sumDur, nil
}
