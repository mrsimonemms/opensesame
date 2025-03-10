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
import { registerAs } from '@nestjs/config';
import { ApiProperty } from '@nestjs/swagger';
import { Type, plainToInstance } from 'class-transformer';
import {
  ArrayMinSize,
  ArrayUnique,
  IsDefined,
  IsNotEmpty,
  IsUrl,
  ValidateNested,
  validateSync,
} from 'class-validator';
import { readFileSync } from 'fs';
import { parse } from 'yaml';

export class Provider {
  @ApiProperty({
    type: 'string',
    description:
      'Provider ID. This is unique and used when identifying the provider programmatically',
    uniqueItems: true,
    example: 'github',
    required: true,
  })
  @IsDefined()
  id: string;

  @ApiProperty({
    type: 'string',
    description: 'Provider name. This is used when creating the login buttons',
    example: 'GitHub',
    required: true,
  })
  @IsDefined()
  name: string;

  @IsDefined()
  @IsUrl({
    protocols: [],
    allow_underscores: true,
    require_tld: false,
    require_protocol: true,
  })
  address: string;
}

export class ProviderConfig {
  @ArrayMinSize(1)
  @ArrayUnique<Provider>((p) => p.id)
  @ValidateNested()
  @Type(() => Provider)
  providers: Provider[];

  @IsDefined()
  @IsNotEmpty()
  protoPath: string;
}

export default registerAs('providers', (): ProviderConfig => {
  const configFilePath = process.env.PROVIDERS_CONFIG_FILE;

  if (!configFilePath) {
    throw new Error('PROVIDERS_CONFIG_FILE must be specified');
  }

  const config = plainToInstance(
    ProviderConfig,
    parse(readFileSync(configFilePath, 'utf-8')),
  );

  // Set the proto path
  config.protoPath = process.env.PROTO_PATH ?? '';

  const err = validateSync(config);
  if (err.length > 0) {
    throw new Error(err.toString());
  }

  return config;
});
