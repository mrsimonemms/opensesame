import { Controller, Inject } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Request } from 'express';
import * as passport from 'passport';
import { Strategy } from 'passport-github2';

import { PassportSDK } from './app.strategy';
import {
  AUTHENTICATION_SERVICE_NAME,
  AuthRequest,
  AuthResponse,
} from './interfaces/authentication/v1/authentication';

@Controller()
export class AppController {
  @Inject('strategy')
  private readonly strategy: Strategy;

  @Inject('passport')
  private readonly passport: PassportSDK;

  @GrpcMethod(AUTHENTICATION_SERVICE_NAME, 'auth')
  auth(data: AuthRequest): AuthResponse {
    console.log('auth');
    console.log(data);

    // console.log(this.strategy);

    this.passport.authenticate(this.strategy, data as unknown as Request);

    return { response: 'oioi' };
  }
}
