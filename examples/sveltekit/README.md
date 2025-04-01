# Svelte Kit

Demo project with [Svelte Kit](https://svelte.dev/docs/kit)

<!-- toc -->

* [Quick start](#quick-start)
* [Authentication journey](#authentication-journey)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

## Quick start

1. Create a [GitHub OAuth app](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app)
   and a [GitLab OAuth app](https://docs.gitlab.com/integration/oauth_provider)
   and store the credentials in a `.envrc` in the root of the project.

   ```sh
   export GITHUB_CLIENT_ID=xxx
   export GITHUB_CLIENT_SECRET=xxx
   export GITLAB_CLIENT_ID=xxx
   export GITLAB_CLIENT_SECRET=xxx
   ```

2. Run the `docker compose up sveltekit` file in the `/examples/sveltekit`
   directory

3. Go to [localhost:9999](http://localhost:9999)

## Authentication journey

In this example, most of the authentication happens on the server. There are no
security reasons why this should happen as no system credentials are required.
This gets deployed to a [Node server](https://svelte.dev/docs/kit/adapter-node),
but could be done anywhere.

This is what could be used in a containerised environment, such as Kubernetes.
The Vite dev server configures a reverse proxy to the Open Sesame server (on `/auth`)
to emulate Ingress controller/API gateway paths.

This journey builds up the request in order. It assumes a certain amount of
knowledge of SvelteKit.

1. [src/hooks.server.ts](./src/hooks.server.ts)

   In SvelteKit, [hooks](https://svelte.dev/docs/kit/hooks) are used to declare
   global functions. In this instance, this checks for a `token` cookie and saves
   it to the [Locals](https://svelte.dev/docs/kit/types#Locals) part of the request.

2. [src/routes/+layout.server.ts](./src/routes/+layout.server.ts)

   Next comes to the global layout. If the `locals` object has a `token` defined,
   it makes a call to get the user from the Open Sesame server. If this call fails,
   the route still continues.

   This maintains the [single-responsibility principle](https://en.wikipedia.org/wiki/Single-responsibility_principle).
   It also shifts the responsibility for deciding what to do with the request to
   the specific groupings, in this case `(protected)` and `(public)`. It's not
   uncommon for there to be more complex logic than what is being demonstrated here.

3. [src/routes/(protected)/+layout.server.ts](./src/routes/(protected)/+layout.server.ts)

   SvelteKit can [group](https://svelte.dev/docs/kit/advanced-routing#Advanced-layouts-(group))
   routes into logical areas. In this example, the groups are used to represent
   sections which can be accessed without valid credentials and those where there
   must be a valid user.

   By the time we get here, the call to the user endpoint has already been made
   so all we need to do is redirect away from here if there is no valid user.

4. [src/routes/(public)/login/+page.svelte](./src/routes/(public)/login/+page.svelte)

   This presents the providers that can be used for login. If you are already logged
   in, it will allow you to add an additional provider to the user.

5. [src/routes/(public)/login/[slug]/+page.server.ts](./src/routes/(public)/login/[slug]/+page.server.ts)

   Arguably, this is not necessary. However as the Open Sesame server requires additional
   query strings (`callback` and `token`) to be sent over, being done by server
   redirection avoids leaking credential recklessly. It also looks nicer to see a
   link without length query strings.

6. [src/routes/(public)/login/callback/+page.server.ts](./src/routes/(public)/login/callback/+page.server.ts)

   This is called when Open Sesame has successfully logged in. This receives the
   token as a [base64 encoded string](https://en.wikipedia.org/wiki/Base64) and
   stores it in the `token` cookie. Once saved, it redirects back to the homepage
   and the whole authentication process validation starts again, but with a token.
