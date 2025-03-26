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

package models

type Pagination[T any] struct {
	Data []T `json:"data"` // Subsection of data

	Count      int32 `json:"count" example:"1"` // Number of records on this page
	Page       int32 `json:"page" example:"2"`  // Current page number
	PerPage    int32 `json:"perPage" example:"25"`
	TotalPages int32 `json:"totalPages" example:"2"`
	Total      int32 `json:"total" example:"26"`
}
