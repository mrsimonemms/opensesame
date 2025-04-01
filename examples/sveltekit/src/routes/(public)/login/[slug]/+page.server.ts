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
import type { PageServerLoad } from './$types';

// We redirect to the Open Sesame URL to avoid displaying tokens in the URL links
export const load: PageServerLoad = async ({ parent, params, url }) => {
  // Check if we're logged in
  const { token, user } = await parent();

  // Get the URL that we'll be redirected to after login
  const callbackURL = new URL(url.toString());
  callbackURL.pathname = '/login/callback';

  let targetURL = `/auth/providers/${params.slug}/login?callback=${encodeURIComponent(callbackURL.toString())}`;

  if (user) {
    // If logged in, add the token to the target URL to link the provider to the existing account
    targetURL += `&token=${token}`;
  }

  return redirect(307, targetURL);
};
