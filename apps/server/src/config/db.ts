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
import { MongooseModuleOptions } from '@nestjs/mongoose';

export default registerAs('db', (): MongooseModuleOptions => {
  return {
    uri: process.env.MONGODB_URL,
    autoCreate: process.env.MONGODB_AUTO_CREATE !== 'false',
    autoIndex: process.env.MONGODB_AUTO_INDEX !== 'false',
    dbName: process.env.MONGODB_DB_NAME,
    minPoolSize: Number(process.env.MONGODB_MIN_POOL_SIZE ?? 5),
    maxPoolSize: Number(process.env.MONGODB_MAX_POOL_SIZE ?? 10),
    retryAttempts: Number(process.env.MONGODB_RETRY_ATTEMPTS ?? 3),
  };
});
