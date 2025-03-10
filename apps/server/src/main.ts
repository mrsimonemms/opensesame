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
import secureSession, {
  SecureSessionPluginOptions,
} from '@fastify/secure-session';
import {
  ClassSerializerInterceptor,
  ConsoleLogger,
  Logger,
  ValidationPipe,
  VersioningType,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { NestFactory, Reflector } from '@nestjs/core';
import {
  FastifyAdapter,
  NestFastifyApplication,
} from '@nestjs/platform-fastify';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

import { AppModule } from './app.module';
import loggerConfig from './config/logger';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
    { logger: new ConsoleLogger(loggerConfig()) },
  );

  const config = app.get(ConfigService);
  const logger = new Logger('bootstrap');

  app
    .enableShutdownHooks()
    .enableVersioning({
      type: VersioningType.URI,
      defaultVersion: '1',
    })
    .useGlobalPipes(
      new ValidationPipe({
        transform: true,
        whitelist: true,
      }),
    )
    .useGlobalInterceptors(new ClassSerializerInterceptor(app.get(Reflector)));

  await app.register(
    secureSession,
    config.getOrThrow<SecureSessionPluginOptions>('session'),
  );

  // Add Swagger documentation
  const docBuilderConfig = new DocumentBuilder()
    .setTitle(process.env.npm_package_name!)
    .setDescription('Instant GitOps deployments')
    .setVersion(process.env.name_package_version ?? 'dev')
    .addBearerAuth()
    .build();

  const documentFactory = SwaggerModule.createDocument(app, docBuilderConfig);

  SwaggerModule.setup('api', app, documentFactory);

  logger.log('Starting server');
  await app.listen(
    config.getOrThrow<number>('server.port'),
    config.getOrThrow('server.host'),
  );
}

bootstrap().catch((err: Error) => {
  /* Unlikely to get to here but a final catchall */
  console.log(err.stack);
  process.exit(1);
});
