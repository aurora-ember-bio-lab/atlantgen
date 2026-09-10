FROM golang:1.22-alpine AS build

RUN apk add --no-cache curl git

WORKDIR /app

# Install Atlas CLI for schema orchestration
RUN curl -sSf https://atlasgo.sh | sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api ./cmd/aura-amber-api

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/api /app/api
COPY --from=build /usr/local/bin/atlas /usr/local/bin/atlas
COPY --from=build /app/atlas ./atlas
ENV DB_URL=postgres://aura:aura@db:5432/aura_amber?sslmode=disable
EXPOSE 8080
CMD ["/app/api"]
