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
// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v2.6.1
//   protoc               unknown
// source: authentication/v1/authentication.proto
/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from '@nestjs/microservices';
import { Observable } from 'rxjs';

export const protobufPackage = 'authentication.v1';

/** KeyRepeatedValue handles definition of repeated values in maps */
export interface KeyRepeatedValue {
  /** Value to use */
  value: string[];
}

/** AuthRequest receives the information about the request */
export interface AuthRequest {
  /** JSON string of the body object */
  body: string;
  /** Headers object */
  headers: { [key: string]: KeyRepeatedValue };
  /** Header method, eg GET, POST, PUT, DELETE etc */
  method: string;
  /** Query object */
  query: { [key: string]: string };
  /** URL, without the domain */
  url: string;
}

export interface AuthRequest_HeadersEntry {
  key: string;
  value: KeyRepeatedValue | undefined;
}

export interface AuthRequest_QueryEntry {
  key: string;
  value: string;
}

/** AuthResponse response for an Auth request */
export interface AuthResponse {
  /** Redirecting */
  redirect?: Redirect | undefined;
  /** Successful call */
  success?: Success | undefined;
}

/** Redirecting the webpage to somewhere else */
export interface Redirect {
  /** URL to redirect to */
  url: string;
  /** HTTP status code */
  status: number;
}

/** Success - valid login information */
export interface Success {
  /** User's login information */
  user: User | undefined;
  /** Info object - stringified JSON */
  info?: string | undefined;
}

/** User - return the user information */
export interface User {
  /** The user ID used by the provider */
  providerId: string;
  /** Any tokens needed to login to the provider */
  tokens: { [key: string]: string };
  /** The user's name according to the provider */
  name?: string | undefined;
  /** The user's username according to the provider */
  username?: string | undefined;
  /** The user's email address according to the provider */
  emailAddress?: string | undefined;
}

export interface User_TokensEntry {
  key: string;
  value: string;
}

export const AUTHENTICATION_V1_PACKAGE_NAME = 'authentication.v1';

/** AuthenticationService handles the individual authentication strategies */

export interface AuthenticationServiceClient {
  /** Handles a new authentication request */

  auth(request: AuthRequest): Observable<AuthResponse>;
}

/** AuthenticationService handles the individual authentication strategies */

export interface AuthenticationServiceController {
  /** Handles a new authentication request */

  auth(request: AuthRequest): Observable<AuthResponse>;
}

export function AuthenticationServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ['auth'];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(
        constructor.prototype,
        method,
      );
      GrpcMethod('AuthenticationService', method)(
        constructor.prototype[method],
        method,
        descriptor,
      );
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(
        constructor.prototype,
        method,
      );
      GrpcStreamMethod('AuthenticationService', method)(
        constructor.prototype[method],
        method,
        descriptor,
      );
    }
  };
}

export const AUTHENTICATION_SERVICE_NAME = 'AuthenticationService';
