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
import { Request } from 'express';
import { Strategy, StrategyFailure } from 'passport';

export interface VerifyCallback {
  (req: Request, done: VerifiedCallback): void;
}

export interface VerifiedCallback {
  (err: Error | null, user?: unknown, info?: unknown): void;
}

export class ProvidersStrategy extends Strategy {
  name = 'providers';

  constructor(private readonly verify: VerifyCallback) {
    super();
  }

  authenticate(req: Request) {
    this.verify(req, (err: Error | null, user?: unknown, info?: unknown) => {
      if (err) {
        return this.error(err);
      }
      if (!user) {
        return this.fail(info as StrategyFailure | string | number);
      }

      this.success(user, info as object);
    });
  }
}
