# Development

<!-- toc -->

* [Dev token](#dev-token)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

This is where development-only resources are kept

## Dev token

If using the standard Docker Compose setup, there is a test user pre-loaded to
the environment to test endpoints. The token to use is:

```txt
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4OTg1NzMxNTksImlhdCI6MTc0MjkwMzEyNiwiaXNzIjoib3BlbnNlc2FtZS5jbG91ZCIsIm5iZiI6MTc0MjkwMzEyNiwic3ViIjoiNTA3ZjFmNzdiY2Y4NmNkNzk5NDM5MDExIn0.Q_JGlCd-QfsFtahdFI5iIFovCh0Q0MN-1B6jXRQnh2A
```

This can either be used as an [HTTP Bearer token](https://datatracker.ietf.org/doc/html/rfc6750)
or with the query string `token=`.
