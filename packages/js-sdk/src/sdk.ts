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
import { STATUS_CODES } from 'http';
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
      const res = await this.exec(strategy, opts);
      if (res) {
        result = res;
        break;
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
    return new Promise<AuthResponse | void>((resolve, reject) => {
      // Create the strategy functions
      const fns = strategy as Strategy & StrategyCreatedStatic;

      /**
       * Internal error while performing authentication.
       *
       * Strategies should call this function when an internal error occurs
       * during the process of performing authentication; for example, if the
       * user directory is not available.
       *
       * @param {Error} err
       */
      fns.error = function (err: Error) {
        // Use reject rather than throwing to avoid this error being repeatedly called
        reject(
          new RpcException({
            code: GrpcStatus.FAILED_PRECONDITION,
            message: err?.message ?? 'Unknown error',
          }),
        );
      };

      /**
       * Fail authentication, with optional `challenge` and `status`, defaulting
       * to 401.
       *
       * Strategies should call this function to fail an authentication attempt.
       *
       * @param {StrategyFailure | String | Number} challenge
       * @param {Number} status
       */
      fns.fail = function (
        challenge?: StrategyFailure | string | number,
        status?: number,
      ) {
        let message = STATUS_CODES[401];

        switch (typeof challenge) {
          case 'string': {
            message = challenge;
            break;
          }
          case 'object': {
            message = challenge.message ?? message;
            break;
          }
        }

        throw new RpcException({
          code:
            status !== 401
              ? GrpcStatus.FAILED_PRECONDITION
              : GrpcStatus.UNAUTHENTICATED,
          message,
        });
      };

      /**
       * Pass without making a success or fail decision.
       *
       * Unlikely to be useful in this application, but exists for PassportJS
       * compatibility
       */
      fns.pass = function () {
        resolve();
      };

      /**
       * Redirect to `url` with optional `status`, defaulting to 302.
       *
       * Strategies should call this function to redirect the user (via their
       * user agent) to a third-party website for authentication.
       *
       * @param {String} url
       * @param {Number} status
       */
      fns.redirect = function (url: string, status: number = 302) {
        resolve({
          redirect: {
            url,
            status,
          },
        });
      };

      /**
       * Authenticate `user`, with optional `info`.
       *
       * Strategies should call this function to successfully authenticate a
       * user.  `user` should be an object supplied by the application after it
       * has been given an opportunity to verify credentials.  `info` is an
       * optional argument containing additional user information.  This is
       * useful for third-party authentication strategies to pass profile
       * details.
       *
       * @param {User} user
       * @param {Object} info
       */
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
