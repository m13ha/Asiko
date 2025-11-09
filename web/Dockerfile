FROM node:20-alpine

WORKDIR /usr/src/app

COPY package*.json ./
COPY tsconfig.json vite.config.ts index.html ./
COPY ./src ./src
COPY ./scripts ./scripts
COPY ./api-client ./api-client

RUN npm install

COPY docker-entrypoint.sh /usr/src/app/docker-entrypoint.sh
RUN chmod +x /usr/src/app/docker-entrypoint.sh

EXPOSE 5173

ENTRYPOINT ["/usr/src/app/docker-entrypoint.sh"]
