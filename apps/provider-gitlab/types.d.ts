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

declare module 'passport-gitlab2' {
  import passport from 'passport';
  import oauth2 from 'passport-oauth2';
  import express from 'express';
  import { OutgoingHttpHeaders } from 'http';

  export interface _StrategyOptionsBase
    extends OAuth2StrategyOptionsWithoutRequiredURLs {
    baseURL: string | undefined;
    clientID: string;
    clientSecret: string;
    callbackURL: string;

    scope?: string[] | undefined;
    userAgent?: string | undefined;
    state?: string | undefined;

    authorizationURL?: string | undefined;
    tokenURL?: string | undefined;
    scopeSeparator?: string | undefined;
    customHeaders?: OutgoingHttpHeaders | undefined;
    userProfileURL?: string | undefined;
    userEmailURL?: string | undefined;
    allRawEmails?: boolean | undefined;
  }

  export interface StrategyOptions extends _StrategyOptionsBase {
    passReqToCallback?: false | undefined;
  }

  export interface StrategyOption extends passport.AuthenticateOptions {
    clientID: string;
    clientSecret: string;
    callbackURL: string;

    scope?: string[] | undefined;
    userAgent?: string | undefined;

    authorizationURL?: string | undefined;
    tokenURL?: string | undefined;
    scopeSeparator?: string | undefined;
    customHeaders?: OutgoingHttpHeaders | undefined;
    userProfileURL?: string | undefined;
    userEmailURL?: string | undefined;
    allRawEmails?: boolean | undefined;
  }

  export class Strategy extends oauth2.Strategy {
    constructor(options: StrategyOptions, verify: oauth2.VerifyFunction);
    constructor(
      options: StrategyOptionsWithRequest,
      verify: oauth2.VerifyFunctionWithRequest,
    );
    userProfile(
      accessToken: string,
      done: (err?: Error | null, profile?: any) => void,
    ): void;

    name: string;
    authenticate(
      req: express.Request,
      options?: passport.AuthenticateOptions,
    ): void;
  }
}
