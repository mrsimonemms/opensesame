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
import { ParamsDictionary } from 'express-serve-static-core';
import { IncomingHttpHeaders } from 'http';

import { AuthRequest } from './interfaces/authentication/v1/authentication';

/**
 * ExpressRequest
 *
 * This converts the gRPC input into an Express-like request object.
 *
 * @link https://github.com/expressjs/express/blob/29d09803c11641d910107793947cefe4c0133358/lib/request.js
 */
export class ExpressRequest<
  P = ParamsDictionary,
  ResBody = never,
  ReqBody = object,
  ReqQuery = AuthRequest['query'],
  Locals extends Record<string, unknown> = Record<string, unknown>,
> implements Partial<Request<P, ResBody, ReqBody, ReqQuery, Locals>>
{
  public readonly body: ReqBody;

  public readonly headers: IncomingHttpHeaders;

  public readonly method: string;

  public readonly query: ReqQuery;

  public readonly url: string;

  constructor(req: AuthRequest) {
    let body: ReqBody;
    try {
      body = JSON.parse(req.body) as ReqBody;
    } catch {
      body = {} as ReqBody;
    }

    this.body = body;
    this.headers = Object.entries(req?.headers ?? {}).reduce(
      (result, [key, { value }]) => {
        result[key] = value.length == 1 ? value[0] : value;
        return result;
      },
      {} as IncomingHttpHeaders,
    );
    this.method = req.method;
    this.query = (req?.query ?? {}) as ReqQuery;
    this.url = req.url;
  }

  get(name: 'set-cookie'): string[] | undefined;
  get(name: string): string | undefined;
  get(name: string): string | string[] | undefined {
    return this.header(name);
  }

  header(name: 'set-cookie'): string[] | undefined;
  header(name: string): string | undefined;
  header(name: unknown): string | string[] | undefined {
    if (!name) {
      throw new TypeError('name argument is required to req.get');
    }

    if (typeof name !== 'string') {
      throw new TypeError('name must be a string to req.get');
    }

    const lc = name.toLowerCase();

    switch (lc) {
      case 'referer':
      case 'referrer':
        return this.headers.referrer || this.headers.referer;
      default:
        return this.headers[lc];
    }
  }
}
