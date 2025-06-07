# build
FROM golang:1.24.2 AS builder

WORKDIR /app

# for caching go mod
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gatewayApp cmd/main.go

# final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gatewayApp .
COPY --from=builder /app/config/config.yaml .

CMD ["./gatewayApp"]
