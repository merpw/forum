# backend/forum

RESTful API server to manage data about users, posts, comments and reactions.

---

Authors: [@maximihajlov](https://github.com/maximihajlov), [@healingdrawing](https://github.com/healingdrawing)
, [@nattikim](https://github.com/nattikim), [@sagarishere](https://github.com/sagarishere)

Solved during studying in Gritlab coding school on Åland, January 2023

---

## [Task description and audit questions](https://github.com/01-edu/public/tree/master/subjects/forum)

## Usage

### Run `go run backend/forum [PARAMS]` to start API server

#### Example: `go run backend/forum` to run on default port 8080

### Params:

- `--port` - port to run API server on (default: 8080)
- `--db` - database file path (default: `./database.db`)

### Environment variables:

- `FORUM_BACKEND_SECRET` - optional, secret header `Internal-Auth` value to access private API endpoints. By default,
  all requests to private endpoints with `Internal-Auth` header will be accepted.

> Note: you can easily generate secret with `uuidgen` command, like this: `FORUM_BACKEND_SECRET=$(uuidgen)`.
> Make sure that this secret is shared with frontend.

- `FORUM_IS_PRIVATE` - optional, default `true`. If `true`, all endpoints will require authentication (except
  `/api/login` and `/api/signup`).

- `FRONTEND_REVALIDATE_URL` - optional, url to revalidate Next.js pages in ISR mode. For
  example, `http://localhost:3000/api/revalidate`
- `FRONTEND_REVALIDATE_TOKEN` - optional, token to revalidate Next.js pages in ISR mode if frontend `/api/` is public

### Testing: `go test backend/forum/... -cover -coverpkg=./...`

### Database migrations

Server supports database migrations. The database will be automatically migrated to the latest version on server start.

If you want to migrate database manually, you can use `cli`, more info
in [backend/migrate](../migrate/README.md) package.

## Endpoints

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/16966820-56131bec-397d-4e40-ad9c-ce1e9b6ec575?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D16966820-56131bec-397d-4e40-ad9c-ce1e9b6ec575%26entityType%3Dcollection%26workspaceId%3D8e6f6f99-c3c2-4738-b609-a958ed3a626a#?env%5BDEV%5D=W3sia2V5IjoiSE9TVCIsInZhbHVlIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwIiwiZW5hYmxlZCI6dHJ1ZSwidHlwZSI6ImRlZmF1bHQifV0=)

- GET `/api/me` - get current user info
- GET `/api/me/posts` - get current user posts
- GET `/api/me/posts/liked` - get current user liked posts
- GET `/api/users/{id}/` - get user info by id
- GET `/api/users/{id}/posts/` - get user posts by id
- GET `/api/posts/` - get all posts
- GET `/api/posts/categories/` - get all categories
- GET `/api/posts/{id}/` - get post by id

- POST `/api/login` - login (get access token by username and password)
- POST `/api/logout` - login (revoke access token)
- POST `/api/signup` - signup (add new user)

- POST `/api/posts/create` - create new post
- POST `/api/posts/{id}/like` - like post by id
- POST `/api/posts/{id}/dislike` - dislike post by id
- POST `/api/posts/{id}/reaction` - get post reaction by id
- POST `/api/posts/{id}/comment/{id}/reaction` - get comment reaction by id
- POST `/api/posts/{id}/comments` - get post comments by id
- POST `/api/posts/{id}/comment` - create new comment
- POST `/api/posts/{id}/comment/{id}/like` - like comment by id
- POST `/api/posts/{id}}/comment/{id}/dislike` - dislike comment by id

[//]: # "TODO: add request body examples"
