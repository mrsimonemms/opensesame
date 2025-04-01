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

import { env } from '$env/dynamic/public';

export type COOKIE_SAME_SITE =
  | true
  | false
  | 'lax'
  | 'strict'
  | 'none'
  | undefined;

export function cookieSameSite(sameSiteVar: string): COOKIE_SAME_SITE {
  let sameSite: COOKIE_SAME_SITE;

  switch (sameSiteVar) {
    case 'true':
      sameSite = true;
      break;
    case 'false':
      sameSite = false;
      break;
    case 'lax':
      sameSite = 'lax';
      break;
    case 'none':
      sameSite = 'none';
      break;
    default:
      sameSite = 'strict';
      break;
  }

  return sameSite;
}

export function tokenCookieOpts(url: URL) {
  return {
    httpOnly: false, // Allow browser to access
    maxAge: Number(env.PUBLIC_COOKIE_EXPIRY ?? 30 * 24 * 60 * 60),
    sameSite: cookieSameSite(env.PUBLIC_COOKIE_SAME_SITE ?? 'lax'),
    path: '/',
    secure: url.protocol === 'https:',
  };
}
