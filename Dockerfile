FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /aura-amber-api ./cmd/aura-amber-api

FROM gcr.io/distroless/static-debian12
COPY --from=build /aura-amber-api /aura-amber-api
EXPOSE 8080
ENTRYPOINT ["/aura-amber-api"]
