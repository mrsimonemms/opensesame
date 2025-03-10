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
import { Controller, Get, Inject, VERSION_NEUTRAL } from '@nestjs/common';
import {
  HealthCheck,
  HealthCheckService,
  MongooseHealthIndicator,
} from '@nestjs/terminus';

@Controller({
  path: 'health',
  version: VERSION_NEUTRAL,
})
export class HealthController {
  @Inject(MongooseHealthIndicator)
  private readonly db: MongooseHealthIndicator;

  @Inject(HealthCheckService)
  private readonly health: HealthCheckService;

  @Get()
  @HealthCheck()
  check() {
    // Allow 1 second before timeout
    const timeout = 1000;

    return this.health.check([
      () =>
        this.db.pingCheck('database', {
          timeout,
        }),
    ]);
  }
}
