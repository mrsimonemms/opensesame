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
import { Controller, Get, Inject } from '@nestjs/common';
import { ApiOperation, ApiResponse } from '@nestjs/swagger';

import { Provider } from '../config/providers';
import { ProvidersService } from './providers.service';

@Controller('providers')
export class ProvidersController {
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
}
