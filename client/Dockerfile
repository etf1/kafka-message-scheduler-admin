FROM node:22-alpine AS builder
WORKDIR /code
COPY package.json .
COPY tsconfig.json .
COPY vite.config.ts .
COPY index.html .
COPY src ./src
COPY public ./public
RUN corepack enable && corepack prepare pnpm@latest --activate
RUN pnpm install --frozen-lockfile && pnpm build
FROM node:22-alpine
RUN addgroup -S app && adduser -S app -G app
RUN mkdir -p /usr/share/ui-dist && chown -R app:app /usr/share/ui-dist
COPY --from=builder /code/build /usr/share/ui-dist
COPY --from=builder /code/server.js /usr/share/ui-dist/server.js
WORKDIR /usr/share/ui-dist
RUN npm i express
USER app
EXPOSE 5000
CMD ["node", "server"]
