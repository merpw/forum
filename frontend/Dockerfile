FROM node:alpine AS base

FROM base AS deps

WORKDIR /app/

COPY package.json package-lock.json ./

RUN npm install

FROM base AS builder

WORKDIR /app/

COPY --from=deps /app/node_modules /app/node_modules
COPY . .

ARG FORUM_BACKEND_PRIVATE_URL

RUN npm run build

FROM base AS runner

WORKDIR /app/

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs --ingroup nodejs

COPY --from=builder /app/public /app/public

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000

VOLUME /app/.next/server/pages

CMD node server.js