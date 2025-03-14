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
import { HttpModule } from '@nestjs/axios';
import { Module } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import {
  ClientGrpcProxy,
  ClientProxyFactory,
  Transport,
} from '@nestjs/microservices';
import { join } from 'path';
import { ProviderConfig } from 'src/config/providers';

import { AUTHENTICATION_V1_PACKAGE_NAME } from '../interfaces/authentication/v1/authentication';
import { PROVIDERS } from './constants';
import { ProvidersController } from './providers.controller';
import { ProvidersService } from './providers.service';

@Module({
  imports: [HttpModule],
  providers: [
    ProvidersService,
    {
      provide: PROVIDERS,
      inject: [ConfigService],
      useFactory: (config: ConfigService) => {
        const { providers, protoPath } =
          config.getOrThrow<ProviderConfig>('providers');

        return providers.reduce((result, provider) => {
          result.set(
            provider.id,
            ClientProxyFactory.create({
              transport: Transport.GRPC,
              options: {
                url: provider.address,
                package: AUTHENTICATION_V1_PACKAGE_NAME,
                protoPath: join(
                  protoPath,
                  'authentication',
                  'v1',
                  'authentication.proto',
                ),
              },
            }),
          );

          return result;
        }, new Map<string, ClientGrpcProxy>());
      },
    },
  ],
  controllers: [ProvidersController],
})
export class ProvidersModule {}
