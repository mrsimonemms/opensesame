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
import { tokenCookieOpts } from '$lib/lib/cookies';

export const load: PageServerLoad = async ({ cookies, url }) => {
  // Get token from query string and decode from base64
  const token = atob(url.searchParams.get('token') ?? '');

  // Set the cookie
  cookies.set('token', token, tokenCookieOpts(url));

  // Go back to homepage and let that handle the redirection
  return redirect(307, '/');
};
