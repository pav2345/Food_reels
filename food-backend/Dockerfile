# =====================================================================
# STAGE 1: Build the Go application binary
# =====================================================================
# Matches your line 1: alpine build environment
FROM golang:1.25-alpine3.22 AS builder

# Set working directory
WORKDIR /app

# Install git (required for fetching certain go modules)
RUN apk update && apk upgrade --no-cache && apk add --no-cache git

# Copy dependency files first for caching benefits
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code (Fixes your incomplete line 16)
COPY . .

# Compile the application into a single portable binary named "server"
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server


FROM alpine:3.22

RUN apk update && apk upgrade --no-cache && apk add --no-cache ca-certificates

WORKDIR /root/

# 1. Copy the compiled application binary from the builder stage
COPY --from=builder /app/server .

# 2. IMPORTANT: Copy your migrations folder so your app can run database migrations at startup
COPY --from=builder /app/migrations ./migrations

# (Optional) If your app loads config files directly from disk, copy them too:
# COPY --from=builder /app/.env.example ./.env

ENV GIN_MODE=release
EXPOSE 8080

ENTRYPOINT ["./server"]