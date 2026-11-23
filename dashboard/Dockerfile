FROM node:20-bookworm-slim AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# The API addresses are included in the build
# NEXT_PUBLIC_API_URL is used by the server side of the dashboard, NEXT_PUBLIC_API_WS_URL by the browser
ARG NEXT_PUBLIC_API_URL=http://127.0.0.1:8081/api/v1
ARG NEXT_PUBLIC_API_WS_URL=ws://127.0.0.1:8081/api/v1/ws
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL NEXT_PUBLIC_API_WS_URL=$NEXT_PUBLIC_API_WS_URL
RUN npm run build && npm prune --omit=dev

FROM node:20-bookworm-slim
WORKDIR /app
ENV NODE_ENV=production
COPY --from=build /app/package.json ./package.json
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/.next ./.next
COPY --from=build /app/public ./public
COPY --from=build /app/next.config.js ./next.config.js
EXPOSE 3000
CMD ["npx", "next", "start", "-p", "3000"]
