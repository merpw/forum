FROM node:alpine AS base

FROM base AS deps

WORKDIR /app/

COPY package.json package-lock.json ./

RUN npm install

FROM deps AS builder

WORKDIR /app/

COPY --from=deps /app/node_modules /app/node_modules
COPY . .

ARG FORUM_BACKEND_PRIVATE_URL

RUN npm run build

FROM builder AS runner

WORKDIR /app/

ENV NODE_ENV=production

COPY --from=builder /app/public /app/public

EXPOSE 3000

ENV PORT 3000

CMD npm run start
