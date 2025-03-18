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
import { Request } from 'express';
import { IncomingHttpHeaders } from 'http';
import { ParsedQs } from 'qs';

import {
  AUTHENTICATION_SERVICE_NAME,
  AuthRequest,
  AuthResponse,
} from './interfaces/authentication/v1/authentication';

@Controller()
export class AppController {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(ConfigService)
  private readonly config: ConfigService;

  @GrpcMethod(AUTHENTICATION_SERVICE_NAME, 'auth')
  auth(data: AuthRequest): AuthResponse {
    const req = {
      body: JSON.parse(data.body) as unknown,
      headers: JSON.parse(data.headers) as IncomingHttpHeaders,
      method: data.method,
      params: JSON.parse(data.params) as unknown,
      query: JSON.parse(data.query) as ParsedQs,
      url: data.url,
    } as Request;

    console.log(req);

    return {};
  }
}
