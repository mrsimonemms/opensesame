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
import { GrpcMethod, RpcException } from '@nestjs/microservices';
import { Strategy } from 'passport';

import { ExpressRequest } from '../express';
import {
  AUTHENTICATION_SERVICE_NAME,
  AuthRequest,
  AuthResponse,
  Route,
  RouteEnabledRequest,
  RouteEnabledResponse,
  UserCreateRequest,
} from '../interfaces/authentication/v1/authentication';
import { SDK } from '../sdk';
import { ROUTES } from './bootstrap';

@Controller()
export class AppController {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject('ROUTES')
  private readonly routes: ROUTES;

  @Inject('STRATEGIES')
  private readonly strategies: Strategy[];

  @GrpcMethod(AUTHENTICATION_SERVICE_NAME, 'auth')
  auth(data: AuthRequest): Promise<AuthResponse> {
    const req = new ExpressRequest(data);

    const passport = new SDK(req);

    return passport.authenticate(this.strategies);
  }

  @GrpcMethod(AUTHENTICATION_SERVICE_NAME, 'routeEnabled')
  routeEnabled(data: RouteEnabledRequest): RouteEnabledResponse {
    return {
      // If it's not defined, treat as disabled
      enabled: this.routes.get(data.route) ?? false,
    };
  }

  @GrpcMethod(AUTHENTICATION_SERVICE_NAME, 'userCreate')
  userCreate(data: UserCreateRequest): UserCreateRequest {
    console.log({ data });

    if (!this.routes.get(Route.ROUTE_USER_CREATE)) {
      // Route not enabled - throw error
      throw new RpcException({
        code: GrpcStatus.NOT_FOUND,
        message: 'Route disabled',
      });
    }

    // @todo(sje): implement user create route
    throw new RpcException({
      code: GrpcStatus.UNIMPLEMENTED,
      message: 'Unimplemented',
    });
  }
}
