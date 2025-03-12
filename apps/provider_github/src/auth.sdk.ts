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
import { Logger } from '@nestjs/common';
import { Request } from 'express';
import {
  AuthenticateOptions,
  Strategy,
  StrategyCreatedStatic,
  StrategyFailure,
} from 'passport';

import { AuthResponse } from './interfaces/authentication/v1/authentication';

export class AuthSDK {
  protected readonly logger = new Logger(this.constructor.name);

  constructor(private readonly req: Request) {}

  authenticate(strategies: Strategy | Strategy[]): AuthResponse;
  authenticate(
    strategies: Strategy | Strategy[],
    opts: AuthenticateOptions,
  ): AuthResponse;
  authenticate(
    strategies: Strategy | Strategy[],
    opts: AuthenticateOptions = {},
  ): AuthResponse {
    if (!Array.isArray(strategies)) {
      strategies = [strategies];
    }

    for (const strategy of strategies) {
      let result: AuthResponse = {};
      this.logger.debug('Attempting strategy login', {
        name: strategy.name ?? 'unnamed',
      });

      // Create the strategy functions
      const fns = strategy as Strategy & StrategyCreatedStatic;

      fns.error = function (err: Error) {
        console.log({ err });
        throw new Error('error not implemented');
      };
      fns.fail = function (
        challenge?: StrategyFailure | string | number,
        status?: number,
      ) {
        console.log({ challenge, status });
        throw new Error('fail not implemented');
      };
      fns.pass = function () {
        throw new Error('pass not implemented');
      };
      fns.redirect = function (url: string, status: number = 302) {
        result = {
          redirect: {
            url,
            status,
          },
        };
      };
      fns.success = function (user: object, info?: object) {
        console.log({ user, info });
        throw new Error('success not implemented');
      };

      this.logger.debug('Calling authentication on strategy');
      fns.authenticate.call(fns, this.req, opts);

      if (result) {
        this.logger.debug('Authenticate call finished');
        return result;
      }
    }

    this.logger.warn('All login strategies have failed');

    console.log({
      //   strategies,
      opts,
    });

    return {};
  }
}
