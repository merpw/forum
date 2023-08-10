# backend/attachments

Image uploading service.

---

## Usage

### Run `go run backend/attachments [PARAMS]` to start API server

_**Example:** `go run .` to run on default port 8082_

### Params:

- `--port` - port to run API server on (default: 8081)
- `--dir` - directory to store uploaded files (default: `./attachments`)

### Environment variables:

- `AUTH_BASE_URL` - optional, default http://localhost:8080 - base url to the auth service

- `FORUM_BACKEND_SECRET` - optional, secret header `Internal-Auth` value to access `/internal/` routes of the Auth
  service
