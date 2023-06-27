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

## Production:

> Note: in production, it's highly recommended to use `FORUM_BACKEND_SECRET` to secure private API endpoints.
>
> This can be done by generating a random string and passing it as an environment variable to docker compose:
>
> `FORUM_BACKEND_SECRET=$(openssl rand -hex 32) docker compose up`

### To run `main` branch version ([forum-ci.mer.pw](https://forum-ci.mer.pw))

```shell
FORUM_BACKEND_SECRET=$(openssl rand -hex 32) docker compose -f docker-compose.yml -f docker-compose.main.yml up
```

### Natively (dev)

requirements: Node.js, Golang, GCC

Check frontend and backend READMEs for more information.

## Environment variables

Frontend and backend have their own environment variables. Check their READMEs for more information.

Here's the list of environment variables used by both:

- `FORUM_BACKEND_SECRET` - optional, secret `Internal-Auth` header value to bypass authentication for private API
  endpoints. If not set, private API endpoints will be available to anyone with `Internal-Auth` header set to any value.

- `FORUM_IS_PRIVATE` - optional, default: `true`.
  Makes all endpoints private by default. If set to `false`, some endpoints will be available to anyone (
  like `/api/posts`).
