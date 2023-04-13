FROM golang:alpine AS builder

# install gcc, required by sqlite3 driver
RUN apk add build-base

WORKDIR /app/

# Download dependencies first to leverage Docker cache
COPY go.* .
RUN go mod download

COPY . .

# Test with cache
RUN --mount=type=cache,target=/root/.cache/go-build go test forum/server/test -cover -coverpkg=./... -coverprofile=./profile.txt

# Build with cache
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /bin/server .

FROM alpine

COPY --from=builder /bin/ /bin/

WORKDIR /app/

VOLUME /app/db

CMD server --db /app/db/database.db

EXPOSE 8080