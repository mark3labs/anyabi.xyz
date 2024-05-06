# syntax=docker/dockerfile:1
# Stage 1: Build the static files
FROM node:22-alpine3.18 as ui-builder
WORKDIR /ui
RUN corepack enable && corepack prepare pnpm@latest --activate
COPY /ui/package.json /ui/pnpm-lock.yaml ./
RUN pnpm i --frozen-lockfile
COPY /ui .
RUN pnpm run build

# Stage 2: Build the Go binary
FROM golang:1.18.2-alpine3.14 AS builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . .
COPY --from=ui-builder /ui/build ./ui/build/
RUN go mod download
RUN CGO_ENABLED=0 go build -o pocketbase .

FROM scratch

WORKDIR /app
COPY --from=builder /app/pocketbase /app/pocketbase
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Start Pocketbase
CMD [ "/app/pocketbase", "serve", "--http=0.0.0.0:5000" ]
