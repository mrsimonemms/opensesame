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
import { status as GrpcStatus } from '@grpc/grpc-js';
import { Controller, Inject, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { GrpcMethod, RpcException } from '@nestjs/microservices';
import {
  Strategy as GitHubStrategy,
  Profile,
  StrategyOptions,
} from 'passport-github2';
import { Strategy as LocalStrategy } from 'passport-local';
import { VerifyCallback } from 'passport-oauth2';

import { ExpressRequest } from './express';
import {
  AUTHENTICATION_SERVICE_NAME,
  AuthRequest,
  AuthResponse,
  Route,
  RouteEnabledRequest,
  RouteEnabledResponse,
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

    const githubStrategy = new GitHubStrategy(
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

    console.log(githubStrategy);

    const localStrategy = new LocalStrategy((username, password, done) => {
      console.log({
        username,
        passport,
      });
      done(null, false);
    });

    const passport = new SDK(req);

    return passport.authenticate([localStrategy]);
  }

  @GrpcMethod(AUTHENTICATION_SERVICE_NAME, 'routeEnabled')
  routeEnabled(data: RouteEnabledRequest): RouteEnabledResponse {
    console.log({ data });
    let enabled: boolean = false;

    switch (data.route) {
      case Route.ROUTE_LOGIN_GET:
        enabled = false;
        break;
      case Route.ROUTE_LOGIN_POST:
        enabled = false;
        break;
      case Route.ROUTE_CALLBACK_GET:
        enabled = true;
        break;
      default:
        throw new RpcException({
          code: GrpcStatus.UNIMPLEMENTED,
        });
    }

    return {
      enabled,
    };
  }
}
