# Build stage
FROM node:slim as builder

WORKDIR /app

COPY . .

RUN npm install

RUN npm run build && docker system prune -f

# Runtime stage
FROM node:slim

COPY --from=builder /app/build /app/build

WORKDIR /app

CMD npm run start

EXPOSE 3000
