import express from "express";
// import passport from "passport";
// import expressSession from "express-session";
import { BasicStrategy } from "passport-http";

const app = express();

// app.use(
//   expressSession({
//     secret: "q1w2e3r4",
//   }),
// );
// app.use(passport.initialize());
// app.use(passport.session());
// passport.use(
//   new BasicStrategy((username, password, done) => {
//     done(null, {
//       username,
//       password,
//     });
//     // console.log({
//     //   username,
//     //   password,
//     //   done,
//     // });
//   }),
// );

app.get("/", (req, res, next) => {
  const req2 = {
    url: req.url,
    body: req.body ?? {},
    headers: req.headers,

    body: { ...req.body },
    headers: { ...req.headers },
    params: { ...req.params },
    query: { ...req.params },
    method: req.method,
    url: req.url,
  };

  const strategy = new BasicStrategy((username, password, done) => {
    done(null, {
      username,
      password,
    });
    // console.log({
    //   username,
    //   password,
    //   done,
    // });
  });

  PassportSDK.authenticate(strategy, { session: false })(req2, res, () => {
    console.log(222);
    res.json(req2.user);
  });
});

app.listen(3000, () => {
  console.log("Listening on 3000");
});

class PassportSDK {
  error(err) {
    console.log({
      err,
    });
  }

  fail(challenge, status) {
    console.log({
      challenge,
      status,
    });
  }

  pass() {
    console.log("pass");
  }

  redirect(url, status = undefined) {
    console.log({ url, status });
  }

  success(user, info = undefined) {
    console.log({ user, info });
  }

  static authenticate(strategy, opts) {
    return (req, res, done) => {
      done(null, { hello: "twatty" });
    };
  }
}
