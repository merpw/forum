FROM golang:alpine AS builder

# install gcc, required by sqlite3 driver
RUN apk add build-base

WORKDIR /app/
COPY . .

# Test with cache
RUN go test ./...

# Build with cache
RUN go build -o /bin/server .

FROM alpine

FROM alpine

COPY --from=builder /bin/ /bin/

WORKDIR /bin/

CMD ./server

EXPOSE 8080