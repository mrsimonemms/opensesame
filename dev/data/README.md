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
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4OTg1NzMxNTksImlhdCI6MTc0MjkwMzEyNiwiaXNzIjoiY2xvdWQtbmF0aXZlLWF1dGgiLCJuYmYiOjE3NDI5MDMxMjYsInN1YiI6IjUwN2YxZjc3YmNmODZjZDc5OTQzOTAxMSJ9.MfozqyuUj7pM8OX9JfYHyRu06JpcrioqBqYh5b8GlYI
```

This can either be used as an [HTTP Bearer token](https://datatracker.ietf.org/doc/html/rfc6750)
or with the query string `token=`.
