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

export interface OrgUserModel {
  userId: string;
  role: string;
  name?: string;
  emailAddress?: string;
  createdDate: Date;
  updatedDate: Date;
}

export interface OrgModel {
  id: string;
  name: string;
  slug: string;
  users: OrgUserModel[];
  createdDate: Date;
  updatedDate: Date;
}
