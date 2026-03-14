FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend
ENV CYPRESS_INSTALL_BINARY=0

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

COPY frontend/ ./
RUN pnpm build

FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN go install github.com/knadh/stuffbin/...@latest
RUN CGO_ENABLED=0 go build -a -ldflags="-s -w" -o canalgov cmd/*.go
RUN "$(go env GOPATH)/bin/stuffbin" -a stuff -in canalgov -out canalgov frontend/dist i18n schema.sql static

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /canalgov

COPY --from=builder /app/canalgov ./canalgov
COPY config.sample.toml ./config.toml

EXPOSE 9000

CMD ["./canalgov"]
