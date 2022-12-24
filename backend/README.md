# FORUM-BACKEND

RESTful API server to manage data about users, posts

---

Authors: [@maximihajlov](https://github.com/maximihajlov), [@healingdrawing](https://github.com/healingdrawing)
, [@nattikim](https://github.com/nattikim)

Solved during studying in Gritlab coding school on Åland, December 2022

---

## [Task description and audit questions](https://github.com/01-edu/public/tree/master/subjects/forum)

## Usage

### Run `go run . --port=[PORT]` to start API server on specified port

#### Example: `go run .` to run on default port 8080

## Endpoints

- `GET /api/posts` - recent posts
- `GET /api/post/{id}` - get post by id
- `POST /api/create` - create new post
- `POST /api/auth/login` - login (get access token by username and password)
- `POST /api/auth/logout` - login (revoke access token)
- `POST /api/auth/signup` - signup (add new user)

[//]: # (TODO: add request body examples)
