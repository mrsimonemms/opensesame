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
import { Server } from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { ReflectionService } from '@grpc/reflection';
import { ConsoleLogger } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { Strategy } from 'passport';
import { join } from 'path';

import {
  AUTHENTICATION_V1_PACKAGE_NAME,
  Route,
} from '../interfaces/authentication/v1/authentication';
import { AppModule } from './app.module';
import loggerConfig from './config/logger';

export type ROUTES = Map<Route, boolean>;

export async function bootstrapPassport(
  strategies: Strategy[],
  routes?: ROUTES,
) {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    AppModule.register(strategies, routes),
    {
      logger: new ConsoleLogger(loggerConfig()),
      transport: Transport.GRPC,
      options: {
        url: process.env.LISTEN_URL ?? '0.0.0.0:3000',
        package: AUTHENTICATION_V1_PACKAGE_NAME,
        protoPath: join(
          process.env.PROTO_PATH ?? '',
          'authentication',
          'v1',
          'authentication.proto',
        ),
        onLoadPackageDefinition: (
          pkg: protoLoader.PackageDefinition,
          server: Pick<Server, 'addService'>,
        ) => {
          new ReflectionService(pkg).addToServer(server);
        },
      },
    },
  );

  await app.listen();
}
