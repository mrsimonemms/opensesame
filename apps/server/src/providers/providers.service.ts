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
import {
  Inject,
  Injectable,
  InternalServerErrorException,
  Logger,
  NotFoundException,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { ClientGrpcProxy } from '@nestjs/microservices';
import { Request as ExpressReq, Response as ExpressRes } from 'express';
import { FastifyReply, FastifyRequest } from 'fastify';
import { Strategy as PassportStrategy } from 'passport';

import { Provider, ProviderConfig } from '../config/providers';
import {
  AUTHENTICATION_SERVICE_NAME,
  AuthenticationServiceClient,
} from '../interfaces/authentication/v1/authentication';
import { PROVIDERS } from './constants';
import { ProvidersStrategy } from './providers.strategy';

@Injectable()
export class ProvidersService {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(ConfigService)
  private readonly configService: ConfigService;

  @Inject(PROVIDERS)
  private readonly providers: Map<string, ClientGrpcProxy>;

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

  findProvider(providerId: string): {
    provider: Provider;
    grpc: ClientGrpcProxy;
  } {
    const provider = this.getProviders().providers.find(
      ({ id }) => id === providerId,
    );

    if (!provider) {
      this.logger.debug('Unknown provider', { providerId });
      throw new NotFoundException(`Unknown provider: ${providerId}`);
    }

    const grpc = this.providers.get(providerId);
    if (!grpc) {
      this.logger.warn('Provider does not have a registered gRPC client');
      throw new NotFoundException(`Unknown provider: ${providerId}`);
    }

    return { provider, grpc };
  }

  getProviders(): ProviderConfig {
    return this.configService.getOrThrow<ProviderConfig>('providers');
  }

  generateStrategy(providerId: string, reply: FastifyReply): PassportStrategy {
    const { grpc } = this.findProvider(providerId);

    const service = grpc.getService<AuthenticationServiceClient>(
      AUTHENTICATION_SERVICE_NAME,
    );

    const strategy = new ProvidersStrategy((req: ExpressReq) => {
      service
        .auth({
          body: JSON.stringify(req.body ?? {}),
          headers: JSON.stringify(req.headers ?? {}),
          method: req.method,
          params: JSON.stringify(req.params ?? {}),
          query: JSON.stringify(req.query ?? {}),
          url: req.url,
        })
        .subscribe((result) => {
          const { redirect } = result;
          if (redirect) {
            this.logger.log('Redirecting');
            reply.redirect(redirect.url, redirect.status ?? 302);
            return;
          }

          this.logger.error('Passport strategy subscribe not handled', {
            result,
          });
          reply.send(new InternalServerErrorException());
        });
    });

    return strategy;
  }
}
