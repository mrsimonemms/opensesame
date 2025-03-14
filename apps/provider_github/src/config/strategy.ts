import { registerAs } from '@nestjs/config';
import { StrategyOptions } from 'passport-github2';

export default registerAs('strategy', (): StrategyOptions => {
  const callbackURL = process.env.CALLBACK_URL;
  if (!callbackURL) {
    throw new Error('CALLBACK_URL is required');
  }

  return {
    clientID: process.env.CLIENT_ID ?? '',
    clientSecret: process.env.CLIENT_SECRET ?? '',
    callbackURL,
    scope: [
      'read:user',
      'user:email',
      ...(process.env.SCOPES ?? '')
        .split(',')
        .map((i) => i.trim())
        .filter((i) => i),
    ],
  };
});
