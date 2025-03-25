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

function data() {
  return {
    meta: {
      createdKey: "createdDate",
      updatedKey: "updatedDate",
    },
    data: [
      {
        _id: "507f1f77bcf86cd799439011",
        emailAddress: "test@test.com",
        name: "Test Testington",
        isActive: true,
        accounts: [
          {
            tokens: {
              accessToken:
                "ooo4+x2VmrcJ2vb7LiQn/wu0XJCA3wlR1WGEVNeCbZngP/rQaoAAu+aweXU6HVEWgc7I",
              refreshToken:
                "rPfJubEh+CzJCzxx+AyCpb7ZhMzw8uxPqTHxaHc16CI0oFsTI/eIAMKchpBjyamRZBD+L6M=",
            },
            providerId: "github",
            providerUserId: "11223344",
            emailAddress: "test@test.com",
            name: "Test Testington",
            username: "testtestington",
            createdDate: new Date(),
            updatedDate: new Date(),
          },
        ],
      },
    ],
  };
}
