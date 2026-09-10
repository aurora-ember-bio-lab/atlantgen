FROM node:22-alpine

WORKDIR /app
COPY package.json tsconfig.json ./
RUN npm install --legacy-peer-deps

COPY src ./src
RUN npm run build

ENV DB_URL=postgres://aura:aura@db:5432/aura_amber?sslmode=disable
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379
ENV WORKER_HTTP_PORT=8090
EXPOSE 8090
CMD ["npm", "run", "start"]
