FROM node:alpine AS builder

WORKDIR /app/

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

RUN npm run build

FROM node:alpine AS runner

WORKDIR /app/

ENV NODE_ENV=production

COPY --from=builder /app/package.json /app/package-lock.json ./

RUN npm ci
COPY --from=builder /app/public /app/public

EXPOSE 3000

ENV PORT=3000

CMD npm run start
