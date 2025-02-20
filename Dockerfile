FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o anyabi .

FROM scratch
WORKDIR /app
COPY --from=builder /app/anyabi /app/anyabi
COPY --from=builder /app/static /app/static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Start Pocketbase
CMD [ "/app/anyabi", "serve", "--http=0.0.0.0:5000" ]
