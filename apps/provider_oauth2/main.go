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

package main

import (
	"github.com/mrsimonemms/cloud-native-auth/apps/provider_oauth2/cmd"
	"github.com/mrsimonemms/cloud-native-auth/packages/authentication"
)

const (
	name        = "oauth2"
	description = "Authenticate Cloud-Native-Auth against OAuth2"
)

func main() {
	authCmd := cmd.New()

	authentication.New(name, description, authCmd).
		Execute()
}
