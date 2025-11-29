# Builder
FROM node:20 AS builder
WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Runtime
FROM node:20-slim
WORKDIR /app

ENV NODE_ENV=production \
    NUXT_HOST=0.0.0.0 \
    NUXT_PORT=3000

COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/package*.json /app/

EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
