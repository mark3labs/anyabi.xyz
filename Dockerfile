FROM golang:1.22.5-alpine AS builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o pocketbase .

FROM scratch
WORKDIR /app
COPY --from=builder /app/pb_public /app/pb_public
COPY --from=builder /app/pocketbase /app/pocketbase
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Start Pocketbase
CMD [ "/app/pocketbase", "serve", "--http=0.0.0.0:5000" ]
