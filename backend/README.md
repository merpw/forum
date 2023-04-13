# FORUM-BACKEND

RESTful API server to manage data about users, posts

---

Authors: [@maximihajlov](https://github.com/maximihajlov), [@healingdrawing](https://github.com/healingdrawing)
, [@nattikim](https://github.com/nattikim), [@sagarishere](https://github.com/sagarishere)

Solved during studying in Gritlab coding school on Åland, January 2023

---

## [Task description and audit questions](https://github.com/01-edu/public/tree/master/subjects/forum)

## Usage

### Run `go run . [PARAMS]` to start API server

#### Example: `go run .` to run on default port 8080

### Params:

- `--port` - port to run API server on (default: 8080)
- `--db` - database file path (default: `./database.db`)

### Environment variables:

- `FRONTEND_REVALIDATE_URL` - optional, url to revalidate Next.js pages in ISR mode. For
  example, `http://localhost:3000/api/revalidate`
- `FRONTEND_REVALIDATE_TOKEN` - optional, token to revalidate Next.js pages in ISR mode if frontend `/api/` is public

### Testing: `go test forum/server/test -cover -coverpkg=./...`

## Endpoints

- GET `/api/me` - get current user info
- GET `/api/me/posts` - get current user posts
- GET `/api/me/posts/liked` - get current user liked posts
- GET `/api/user/{id}/` - get user info by id
- GET `/api/user/{id}/posts/` - get user posts by id
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
