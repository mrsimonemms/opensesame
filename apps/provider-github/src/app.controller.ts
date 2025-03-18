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
import { Controller, Inject, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { GrpcMethod } from '@nestjs/microservices';
import {
  Strategy as GitHubStrategy,
  Profile,
  StrategyOptions,
} from 'passport-github2';
import { VerifyCallback } from 'passport-oauth2';

import { ExpressRequest } from './express';
import {
  AUTHENTICATION_SERVICE_NAME,
  AuthRequest,
  AuthResponse,
  User,
} from './interfaces/authentication/v1/authentication';
import { SDK } from './sdk';

@Controller()
export class AppController {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(ConfigService)
  private readonly config: ConfigService;

  @GrpcMethod(AUTHENTICATION_SERVICE_NAME, 'auth')
  async auth(data: AuthRequest): Promise<AuthResponse> {
    const req = new ExpressRequest(data);

    const strategy = new GitHubStrategy(
      this.config.getOrThrow<StrategyOptions>('strategy'),
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

    const passport = new SDK(req);

    return passport.authenticate([strategy]);
  }
}
