# =========================
# 1️⃣ Build Stage
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (sometimes needed)
RUN apk add --no-cache git

# Copy dependency files first (cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# =========================
# 2️⃣ Runtime Stage
# =========================
FROM alpine:latest

WORKDIR /app

# Copy ONLY the compiled binary
COPY --from=builder /app/app .

EXPOSE 8000

CMD ["./app"]
