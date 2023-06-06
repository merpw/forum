# forum

---

Authors: [@maximihajlov](https://github.com/maximihajlov), [@healingdrawing](https://github.com/healingdrawing)
, [@nattikim](https://github.com/nattikim), [@sagarishere](https://github.com/sagarishere)

Solved during studying in Gritlab coding school on Åland, January 2023

---

## [Task description and audit questions](https://github.com/01-edu/public/tree/master/subjects/forum)

---

## Demo [forum.mer.pw](https://forum.mer.pw/)

---

## How to run?

## Docker compose: `docker compose up`

> Note: in production, it's highly recommended to use `FORUM_BACKEND_SECRET` to secure private API endpoints.
>
> You can generate it using the following command:

### Docker production: `FORUM_BACKEND_SECRET=$(uuidgen) docker compose up`

### Natively (dev)

requirements: Node.js, Golang, GCC

#### Commands:

#### Backend: `cd backend && go run backend/forum`

#### Frontend: `cd frontend && npm run dev`

## Environment variables

Frontend and backend have their own environment variables. Check their READMEs for more information.

Here's the list of environment variables used by both:

- `FORUM_BACKEND_SECRET` - optional, secret `Internal-Auth` header value to bypass authentication for private API
  endpoints. If not set, private API endpoints will be available to anyone with `Internal-Auth` header set to any value.

- `FORUM_IS_PRIVATE` - optional, default: `true`.
  Makes all endpoints private by default. If set to `false`, some endpoints will be available to anyone (
  like `/api/posts`).
