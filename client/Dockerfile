FROM node:16.0-alpine as builder

WORKDIR /code

COPY package.json .
COPY yarn.lock .
COPY tsconfig.json .
COPY server.js .
COPY src ./src
COPY public ./public
RUN yarn && yarn build


FROM node:16.0-alpine

RUN mkdir -p /usr/share/ui-dist
COPY --from=builder /code/build /usr/share/ui-dist
COPY --from=builder /code/server.js /usr/share/ui-dist/server.js
WORKDIR /usr/share/ui-dist
RUN npm i express
CMD ["node", "server"]
