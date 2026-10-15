# ----------------------------------------------------
# Stage 1: Build the binary
# ----------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cài đặt git và certificates nếu cần
RUN apk add --no-cache git ca-certificates tzdata

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code
COPY . .

# Build binary tĩnh không phụ thuộc CGO (Static Binary)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/api

# ----------------------------------------------------
# Stage 2: Runtime image tối ưu và an toàn
# ----------------------------------------------------
FROM alpine:3.21

# Cài đặt tzdata, ca-certificates và curl phục vụ healthcheck
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# Copy binary từ stage builder
COPY --from=builder /app/bin/server /app/server

# Expose cổng ứng dụng
EXPOSE 8080

# Healthcheck định kỳ kiểm tra container
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Chạy server
CMD ["/app/server"]
