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
  HealthCheckResult,
  HealthCheckService,
  HealthIndicatorResult,
  MongooseHealthIndicator,
  TerminusModule,
} from '@nestjs/terminus';
import { Test, TestingModule } from '@nestjs/testing';

import { HealthController } from './health.controller';

describe('HealthController', () => {
  let controller: HealthController;
  let service: Partial<HealthCheckService>;
  let dbCheck: Partial<MongooseHealthIndicator>;

  beforeEach(async () => {
    service = {
      check: jest.fn(),
    };

    dbCheck = {
      pingCheck: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      imports: [TerminusModule],
      controllers: [HealthController],
      providers: [
        {
          provide: HealthCheckService,
          useValue: service,
        },
        {
          provide: MongooseHealthIndicator,
          useValue: dbCheck,
        },
      ],
    }).compile();

    controller = module.get<HealthController>(HealthController);
  });

  it('should return a successful result', async () => {
    const serviceResult: HealthCheckResult = {
      status: 'ok',
      details: {},
    };
    const dbResult: HealthIndicatorResult = {
      database: {
        status: 'up',
      },
    };
    const serviceMock = jest
      .spyOn(service, 'check')
      .mockResolvedValue(serviceResult);

    jest.spyOn(dbCheck, 'pingCheck').mockResolvedValue(dbResult);

    await expect(controller.check()).resolves.toBe(serviceResult);

    expect(service.check).toHaveBeenCalled();

    await expect(serviceMock.mock.calls[0][0][0]()).resolves.toBe(dbResult);
  });
});
