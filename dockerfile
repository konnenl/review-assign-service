FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]