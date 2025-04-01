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

import type { UserModel } from '@opensesame-cloud/js-sdk/dist/models/user';
import type { LayoutServerLoad } from './$types';

// This searches for a token in the local request (set by /src/hooks.server.ts)
// and verifies that against the Open Sesame server. If it's valid, it loads
// the user data.
export const load: LayoutServerLoad = async ({ fetch, locals }) => {
  // Get the token from the local request
  const { token } = locals;

  // Make a call to the Open Sesame server
  const res = await fetch('/auth/user', {
    headers: {
      authorization: `Bearer ${token}`,
    },
  });

  // If the user is valid, get the data
  let user: UserModel | undefined;
  if (res.ok) {
    user = (await res.json()).user;
  }

  return {
    token,
    user,
  };
};
