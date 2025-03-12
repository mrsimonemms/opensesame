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
import { ConsoleLogger } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';

import { AppModule } from './app.module';
import loggerConfig from './config/logger';

async function bootstrap() {
  const protoRoot = join(__dirname, '..', '..', '..', 'proto');

  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    AppModule,
    {
      logger: new ConsoleLogger(loggerConfig()),
      transport: Transport.GRPC,
      options: {
        package: 'authentication',
        protoPath: join(
          protoRoot,
          'authentication',
          'v1',
          'authentication.proto',
        ),
      },
    },
  );

  await app.listen();
}

bootstrap().catch((err: Error) => {
  /* Unlikely to get to here but a final catchall */
  console.log(err.stack);
  process.exit(1);
});
