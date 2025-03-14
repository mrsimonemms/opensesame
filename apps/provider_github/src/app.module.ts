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
import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { Strategy, StrategyOptions } from 'passport-github2';
import { BasicStrategy } from 'passport-http';
import { VerifyCallback } from 'passport-oauth2';

import { AppController } from './app.controller';
import { PassportSDK } from './app.strategy';
import config from './config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: config,
    }),
  ],
  controllers: [AppController],
  providers: [
    {
      provide: 'passport',
      useClass: PassportSDK,
    },
    {
      provide: 'strategy',
      useFactory: () => {
        return new BasicStrategy(
          (username: string, password: string, done: unknown) => {
            console.log({
              username,
              password,
              done,
            });
          },
        );
      },
    },
    // {
    //   provide: 'strategy',
    //   inject: [ConfigService],
    //   useFactory: (config: ConfigService) => {
    //     return new Strategy(
    //       config.getOrThrow<StrategyOptions>('strategy'),
    //       (
    //         accessToken: string,
    //         refreshToken: string,
    //         profile: unknown,
    //         done: VerifyCallback,
    //       ) => {
    //         console.log({
    //           accessToken,
    //           refreshToken,
    //           profile,
    //           done,
    //         });
    //       },
    //     );
    //   },
    // },
  ],
})
export class AppModule {}
