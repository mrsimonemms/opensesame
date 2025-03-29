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
import { Route, User, bootstrapPassport } from '@opensesame-cloud/js-sdk';
import { Profile } from 'passport';
import { Strategy, StrategyOptions } from 'passport-github2';
import { VerifyCallback } from 'passport-oauth2';

const callbackURL = process.env.CALLBACK_URL;
if (!callbackURL) {
  throw new Error('CALLBACK_URL is required');
}

const config: StrategyOptions = {
  clientID: process.env.CLIENT_ID ?? '',
  clientSecret: process.env.CLIENT_SECRET ?? '',
  callbackURL,
  scope: [
    'read:user',
    'user:email',
    ...(process.env.SCOPES ?? '')
      .split(',')
      .map((i) => i.trim())
      .filter((i) => i),
  ],
};

const strategy = new Strategy(
  config,
  (
    accessToken: string,
    refreshToken: string,
    profile: Profile,
    done: VerifyCallback,
  ) => {
    const user: User = {
      providerId: profile.id,
      tokens: {
        accessToken,
        refreshToken,
      },
      name: profile.displayName,
      emailAddress: profile.emails?.[0]?.value,
      username: profile.username,
    };

    done(null, user);
  },
);

// Configure the routes to use
const routes = new Map<Route, boolean>([
  [Route.ROUTE_LOGIN_GET, true],
  [Route.ROUTE_CALLBACK_GET, true],
]);

// GO GO GO!!!
bootstrapPassport([strategy], routes).catch((err: Error) => {
  console.log(err.stack);
  process.exit(1);
});
