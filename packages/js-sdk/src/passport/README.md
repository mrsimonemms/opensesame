# Passport

This bootstraps a standalone provider generated from [PassportJS](https://www.passportjs.org/)
strategies.

```ts
import { User, bootstrapPassport } from '@cloud-native-auth/js-sdk';
import { Strategy as GitHubStrategy } from 'passport-github2';

// Create the strategy as-per the PassportJS docs
// @link https://www.passportjs.org/packages/passport-github2/
const strategy = new GitHubStrategy()

bootstrapPassport([strategy]).catch((err: Error) => {
  console.log(err.stack);
  process.exit(1);
});
```
