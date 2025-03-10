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
import { SecureSessionPluginOptions } from '@fastify/secure-session';
import { registerAs } from '@nestjs/config';

function stringToBool(
  str?: string,
  defaultValue?: boolean,
): boolean | undefined {
  if (str === 'true') {
    return true;
  } else if (str === 'false') {
    return false;
  }
  return defaultValue;
}

export default registerAs('session', (): SecureSessionPluginOptions => {
  let sameSite: true | false | 'lax' | 'strict' | 'none' | undefined;
  const sameSiteVar = process.env.SESSION_SAME_SITE;

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

  return {
    cookie: {
      domain: process.env.SESSION_COOKIE_DOMAIN,
      httpOnly: stringToBool(process.env.SESSION_COOKIE_HTTPONLY, true),
      maxAge: Number(process.env.SESSION_COOKIE_MAX_AGE ?? 30 * 24 * 60 * 60),
      path: process.env.SESSION_COOKIE_PATH ?? '/',
      sameSite,
      secure: stringToBool(process.env.SESSION_COOKIE_SECURE),
      signed: stringToBool(process.env.SESSION_COOKIE_SIGNED),
    },
    cookieName: process.env.SESSION_COOKIE_NAME ?? 'session',
    expiry: Number(process.env.SESSION_EXPIRY_SECONDS ?? 30 * 60 * 60), // Default to a short length as only used for linking OAuth sessions
    secret: process.env.SESSION_SECRET ?? '',
    salt: process.env.SESSION_SALT ?? '',
  };
});
