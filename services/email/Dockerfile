FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o build/email services/email/cmd/main.go

FROM alpine
COPY --from=builder /app/build/email /email
EXPOSE 8081
ENTRYPOINT ["./email"]

