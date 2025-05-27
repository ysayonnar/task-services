# build
FROM golang:1.24.2 AS builder

WORKDIR /app

# for caching go mod
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ssoApp cmd/main.go

# final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ssoApp .
COPY --from=builder /app/config/config.yaml .

CMD ["./ssoApp"]
