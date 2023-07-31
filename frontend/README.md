# frontend

Next.js based frontend for [real-time-forum](https://github.com/01-edu/public/tree/master/subjects/real-time-forum) task

---

### Environment variables

- `FORUM_IS_PRIVATE` - optional, default: `true`.
  Enables authentication middleware, redirects to login page if a user is not logged in.

- `FORUM_BACKEND_PRIVATE_URL` - default for `dev`: `http://localhost:8080`.
  URL to [backend](../backend) instance to use for server-side rendering. Optional for building (to pre-render pages),
  but required for running.
- `DEV_FORUM_BACKEND_REWRITE_URL` - optional, default for `dev`: `FORUM_BACKEND_PRIVATE_URL`.
  URL to rewrite all `/api/` requests to backend using Next.js. It should not be used in production, use reverse proxy
  instead.
- `FRONTEND_REVALIDATE_TOKEN` - optional, should be set if `/api/revalidate` endpoint is going to be public. By default,
  it is not public and available only in private network.

- `OPENAI_API_KEY` - optional. API key for [OpenAI](https://openai.com/) API. If set, it will be used
  to generate post descriptions.

## Usage

You need Node.js installed to run frontend separately

### Run `npm install`

To get project dependencies

### Run `npm run dev`

To start development server

### Run `npm run build`

To build optimized production build.

### Run `npm run start`

To start Next.js server

[//]: # "TODO: add comment about data fetching and export mode config"
