FROM golang:alpine AS builder

# for gcc
RUN apk add build-base

WORKDIR /app/
COPY . .

RUN go build -o /bin/server .

FROM alpine

FROM alpine

COPY --from=builder /bin/ /bin/

WORKDIR /bin/

CMD ./server

EXPOSE 8080