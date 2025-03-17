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
  Controller,
  Get,
  HttpStatus,
  Inject,
  InternalServerErrorException,
  Logger,
  Param,
  Req,
  Res,
  UseGuards,
} from '@nestjs/common';
import { AuthGuard } from '@nestjs/passport';
import { ApiOperation, ApiParam, ApiResponse } from '@nestjs/swagger';
import { Handler } from 'express';
import { FastifyReply, FastifyRequest } from 'fastify';
import * as passport from 'passport';

import { Provider } from '../config/providers';
import { ProvidersService } from './providers.service';

@Controller('providers')
export class ProvidersController {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(ProvidersService)
  private readonly service: ProvidersService;

  @Get('/')
  @ApiOperation({
    summary: 'Return list of available providers',
  })
  @ApiResponse({
    status: 200,
    description: 'Return list of available providers',
    type: Provider,
    isArray: true,
  })
  list() {
    return this.service
      .getProviders()
      .providers.map(({ id, name }) => ({ id, name }))
      .sort();
  }

  @Get('/:providerId/login')
  @UseGuards(AuthGuard(['anonymous']))
  @ApiOperation({
    summary:
      'Dispatch to authentication provider login. If user provided, the login will be added to the user',
  })
  @ApiParam({
    name: 'providerId',
    description: 'Provider ID',
    example: 'github',
  })
  @ApiResponse({
    status: HttpStatus.OK,
    description: 'User authenticated',
  })
  @ApiResponse({
    status: HttpStatus.FOUND,
    description: 'Dispatch to authentication provider login',
  })
  @ApiResponse({
    status: HttpStatus.UNAUTHORIZED,
    description: 'User not authenticated',
  })
  @ApiResponse({
    status: HttpStatus.NOT_FOUND,
    description: 'Provider not found',
  })
  login(
    @Param('providerId') providerId: string,
    @Req() req: FastifyRequest,
    @Res() res: FastifyReply,
  ) {
    const strategy = this.service.generateStrategy(providerId, res);

    const handler = passport.authenticate(
      strategy,
      (err: Error, user?: unknown) => {
        if (err) {
          res.send(err);
          return;
        }
        res.send(user);
      },
    ) as Handler;

    return handler(
      ...this.service.fastifyToExpress(req, res),
      (err?: Error | 'router' | 'route') => {
        // If we've gotten here, something has gone very wrong
        this.logger.error(
          'Provider middleware nextfunction has been triggered',
          {
            err,
          },
        );

        throw new InternalServerErrorException(err);
      },
    );
  }

  @Get('/:providerId/login/callback')
  @ApiParam({
    name: 'providerId',
    description: 'Provider ID',
    example: 'github',
  })
  loginCallback(
    @Param('providerId') providerId: string,
    @Req() req: FastifyRequest,
    @Res() res: FastifyReply,
  ) {
    const strategy = this.service.generateStrategy(providerId, res);

    const handler = passport.authenticate(
      strategy,
      (err: Error, user?: unknown) => {
        if (err) {
          res.send(err);
          return;
        }
        res.send(user);
      },
    ) as Handler;

    return handler(
      ...this.service.fastifyToExpress(req, res),
      (err?: Error | 'router' | 'route') => {
        // If we've gotten here, something has gone very wrong
        this.logger.error(
          'Provider middleware nextfunction has been triggered',
          {
            err,
          },
        );

        throw new InternalServerErrorException(err);
      },
    );
  }
}
