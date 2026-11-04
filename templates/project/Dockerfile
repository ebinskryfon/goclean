# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.23-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

# ---- runtime stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -u 10001 app
WORKDIR /app

COPY --from=build /out/api /out/migrate /app/
COPY migrations /app/migrations

USER app
ENV APP_ENV=production APP_HOST=0.0.0.0 APP_PORT=8080
EXPOSE 8080
CMD ["/app/api"]
