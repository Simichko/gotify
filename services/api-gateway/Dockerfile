FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o build/gateway services/api-gateway/cmd/main.go

FROM alpine
COPY --from=builder /app/build/gateway /gateway
ENTRYPOINT ["./gateway"]

