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
import { RpcException } from '@nestjs/microservices';
import {
  AuthenticateOptions,
  Strategy,
  StrategyCreatedStatic,
  StrategyFailure,
} from 'passport';

import { ExpressRequest } from './express';
import {
  AuthResponse,
  User,
} from './interfaces/authentication/v1/authentication';

export type AuthCallback = () => void;

export class SDK {
  constructor(private readonly req: ExpressRequest) {}

  /**
   * Authenticate
   *
   * Implements the PassportJS authentication pattern
   *
   * @link https://github.com/jaredhanson/passport/blob/master/lib/middleware/authenticate.js#L70
   */
  async authenticate(
    strategies: Strategy[],
    opts: AuthenticateOptions = {},
  ): Promise<AuthResponse> {
    let result: AuthResponse | undefined;
    for (const strategy of strategies) {
      console.log(strategy.name);
      console.log('strategy invoking');
      try {
        const r = await this.exec(strategy, opts);
        if (r) {
          result = r;
          break;
        }
      } catch (err: unknown) {
        throw new RpcException({
          code: GrpcStatus.NOT_FOUND,
          message: (err as Error).message ?? '',
        });
        console.log({ err });
      }
    }

    if (!result) {
      // If we get here, everything has failed
      throw new Error('All strategies have failed');
    }

    return result ?? {};
  }

  private async exec(
    strategy: Strategy,
    opts: AuthenticateOptions,
  ): Promise<AuthResponse | void> {
    return new Promise<AuthResponse>((resolve, reject) => {
      // Create the strategy functions
      const fns = strategy as Strategy & StrategyCreatedStatic;

      // @todo(sje): do something with the reject
      console.log({ reject });

      fns.error = function (err: Error) {
        console.log({ err });
        throw new RpcException({
          code: GrpcStatus.UNIMPLEMENTED,
          message: 'error not implemented',
        });
      };
      fns.fail = function (
        challenge?: StrategyFailure | string | number,
        status?: number,
      ) {
        console.log({ challenge, status });
        throw new RpcException({
          code: GrpcStatus.UNIMPLEMENTED,
          message: 'fail not implemented',
        });
      };
      fns.pass = function () {
        throw new RpcException({
          code: GrpcStatus.UNIMPLEMENTED,
          message: 'pass not implemented',
        });
      };
      fns.redirect = function (url: string, status: number = 302) {
        resolve({
          redirect: {
            url,
            status,
          },
        });
      };
      fns.success = function (user: User, info?: object) {
        // Search and remove any undefined tokens
        user.tokens = Object.fromEntries(
          Object.entries(user.tokens).filter(
            ([, value]: [string, string?]): boolean => !!value,
          ),
        );

        let infoStr: string | undefined;
        if (info) {
          infoStr = JSON.stringify(info);
        }

        resolve({
          success: {
            user,
            info: infoStr,
          },
        });
      };

      fns.authenticate.call(fns, this.req, opts);
    });
  }
}
