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
import { HttpService } from '@nestjs/axios';
import { Inject, Injectable, Logger, NotFoundException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Request as ExpressReq, Response as ExpressRes } from 'express';
import { FastifyReply, FastifyRequest } from 'fastify';
import { Strategy as PassportStrategy } from 'passport';

import { Provider, ProviderConfig } from '../config/providers';
import { ProvidersStrategy, VerifiedCallback } from './providers.strategy';

@Injectable()
export class ProvidersService {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(ConfigService)
  private readonly configService: ConfigService;

  @Inject(HttpService)
  private readonly httpService: HttpService;

  fastifyToExpress(
    req: FastifyRequest,
    res: FastifyReply,
  ): [ExpressReq, ExpressRes] {
    const expressRes = res as unknown as ExpressRes;

    expressRes.setHeader = function (
      this: FastifyReply,
      name: string,
      value: number | string | readonly string[],
    ): ExpressRes {
      this.raw.setHeader(name, value);
      return this as unknown as ExpressRes;
    };

    expressRes.end = function (this: FastifyReply): ExpressRes {
      // Use send so the session cookie onSend hook is triggered
      this.send();
      return this as unknown as ExpressRes;
    };

    return [
      // Request has no changes requires
      req as unknown as ExpressReq,
      expressRes,
    ];
  }

  findProvider(providerId: string): Provider {
    const provider = this.getProviders().providers.find(
      ({ id }) => id === providerId,
    );

    if (!provider) {
      this.logger.debug('Unknown provider', { providerId });
      throw new NotFoundException(`Unknown provider: ${providerId}`);
    }

    return provider;
  }

  getProviders(): ProviderConfig {
    return this.configService.getOrThrow<ProviderConfig>('providers');
  }

  generateStrategy(providerId: string): PassportStrategy {
    const provider = this.findProvider(providerId);

    const strategy = new ProvidersStrategy(
      (req: ExpressReq, done: VerifiedCallback) => {
        done(new NotFoundException(provider.id));
      },
    );

    return strategy;
  }
}
