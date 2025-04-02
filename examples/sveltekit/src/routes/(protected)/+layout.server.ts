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

import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { tokenCookieOpts } from '$lib/lib/cookies';

// All routes in here require a valid user. This hook loads the
// user from the local request and redirects to the login page
// if there is no user present.
export const load: LayoutServerLoad = async ({ cookies, parent, url }) => {
  const { token, user } = await parent();

  if (!user) {
    // Delete the bad cookie
    cookies.delete('token', tokenCookieOpts(url));

    // Redirect to the login page
    return redirect(307, '/login');
  }

  return {
    token,
    user,
  };
};
