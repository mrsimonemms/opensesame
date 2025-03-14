import { Request } from 'express';
import { stat } from 'fs';
import * as passport from 'passport';
import { Strategy } from 'passport-strategy';

// passport.authenticate()

// export class PassportSDK implements passport.SessionStrategy {
//   // authenticate(strategy: passport.Strategy) {
//   //   console.log(strategy);
//   // }
// }

// passport.authenticate(
//   new Strategy(
//     {
//       clientID: '',
//       clientSecret: '',
//       callbackURL: '',
//     },
//     () => {},
//   ),
// );

export class S {
  constructor(private readonly strategy: passport.Strategy) {}

  name?: string = 'example';

  authenticate(req: Request, options?: unknown) {
    console.log(this.strategy.authenticate.call(this, req, options));
  }

  error(err: unknown): void {
    console.log({
      err,
    });
  }

  fail(
    challenge?: passport.StrategyFailure | string | number,
    status?: number,
  ): void {
    console.log({
      challenge,
      status,
    });
  }

  pass(): void {
    console.log('pass');
  }

  redirect(url: string, status?: number): void {
    console.log({ url, status });
  }

  success(user: Express.User, info?: object): void {
    console.log({ user, info });
  }
}

export class PassportSDK {
  authenticate(
    strategy: passport.Strategy,
    req: Request,
    options: passport.AuthenticateOptions = {},
  ): void {
    const s = new S(strategy);

    s.authenticate(req, options);

    // console.log(strategy.authenticate.call(new S(), {} as Request, options));

    // strategy.authenticate.apply({});

    // strategy.authenticate();
  }
}
