# FORUM-FRONTEND

Next.js based frontend for [FORUM](https://github.com/01-edu/public/tree/master/subjects/forum) task

---

Authors: [@maximihajlov](https://github.com/maximihajlov), [@healingdrawing](https://github.com/healingdrawing)
, [@nattikim](https://github.com/nattikim), [@sagarishere](https://github.com/sagarishere)

Solved during studying in Gritlab coding school on Åland, January 2023

---

### Environment variables

- `FORUM_BACKEND_PRIVATE_URL` - default for `dev`: `http://localhost:8080`.
  URL to [backend](../backend) instance to use for server-side rendering. Optional for building (to pre-render pages),
  but required for running.
- `DEV_FORUM_BACKEND_REWRITE_URL` - optional, default for `dev`: `FORUM_BACKEND_PRIVATE_URL`.
  URL to rewrite all `/api/` requests to backend using Next.js. It should not be used in production, use reverse proxy
  instead.
- `FRONTEND_REVALIDATE_TOKEN` - optional, should be set if `/api/revalidate` endpoint is going to be public. By default,
  it is not public and available only in private network.

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
