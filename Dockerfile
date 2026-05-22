# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

EXPOSE 8081

CMD ["./main"]