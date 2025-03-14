import express from "express";
import passport from "passport";
import expressSession from "express-session";
import { BasicStrategy } from "passport-http";

const app = express();

app.use(
  expressSession({
    secret: "q1w2e3r4",
  }),
);
app.use(passport.initialize());
app.use(passport.session());
passport.use(
  new BasicStrategy((username, password, done) => {
    done(null, {
      username,
      password,
    });
    // console.log({
    //   username,
    //   password,
    //   done,
    // });
  }),
);

app.get(
  "/",
  (req, res, next) => {
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
    passport.authenticate("basic", { session: false })(req2, res, next);
  },
  (req, res) => {
    res.json(req.user);
  },
);

app.listen(3000, () => {
  console.log("Listening on 3000");
});
