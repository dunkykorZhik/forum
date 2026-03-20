# ---------- BUILD STAGE ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for modules sometimes)
RUN apk add --no-cache git

# Copy go mod first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all code
COPY . .

# Build binary
RUN go build -o server ./cmd


# ---------- RUN STAGE ----------
FROM alpine:latest

WORKDIR /app

# Install certs (important for HTTPS if any)
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/server .

# Copy static + templates
COPY --from=builder /app/internal/web-storage ./internal/web-storage

# (optional) create db folder
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]