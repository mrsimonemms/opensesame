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
import { DynamicModule, Logger, Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { Strategy } from 'passport';

import { Route } from '../interfaces/authentication/v1/authentication';
import { AppController } from './app.controller';
import { ROUTES } from './bootstrap';
import config from './config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: config,
    }),
  ],
  controllers: [AppController],
})
export class AppModule {
  protected static readonly logger = new Logger('AppModule');

  static register(strategies: Strategy[], routes?: ROUTES): DynamicModule {
    return {
      module: AppModule,
      providers: [
        {
          provide: 'STRATEGIES',
          useFactory: (): Strategy[] => {
            if (strategies.length === 0) {
              throw new Error('At least one strategy required');
            }
            return strategies;
          },
        },
        {
          provide: 'ROUTES',
          useFactory: (): ROUTES => {
            if (routes) {
              this.logger.debug('Routes', { routes });
              return routes;
            }

            // By default, enable every route
            this.logger.debug('All routes enabled');
            return new Map<Route, boolean>([
              [Route.ROUTE_LOGIN_GET, true],
              [Route.ROUTE_LOGIN_POST, true],
              [Route.ROUTE_CALLBACK_GET, true],
            ]);
          },
        },
      ],
    };
  }
}
