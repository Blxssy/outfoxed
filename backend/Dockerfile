FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/app/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main ./main

EXPOSE 8080

CMD ["./main"]