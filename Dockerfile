# --- STAGE 1: Build (Koki Masak) ---
FROM golang:alpine AS builder

WORKDIR /app

# Download dependency dulu (biar cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary (output file: main)
# Perhatikan path: ./cmd/main.go
RUN go build -o main ./cmd/main.go

# --- STAGE 2: Run (Penyajian Ringan) ---
FROM alpine:latest

WORKDIR /root/

# Copy binary hasil build
COPY --from=builder /app/main .

# PENTING: Copy folder HTML/Templates!
COPY --from=builder /app/web ./web

# Expose port (formalitas)
EXPOSE 8080

# Gas nyalain!
CMD ["./main"]