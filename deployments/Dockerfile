# Multi-stage Docker build for kbFood with React frontend

# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci --legacy-peer-deps

# Copy frontend source
COPY frontend/ ./

# Build the frontend
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.23-alpine AS backend-builder

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/server

# Stage 3: Final Runtime Image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Set timezone
ENV TZ=Asia/Shanghai

# Create app directory and user
RUN mkdir /app && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser -h /app appuser

WORKDIR /app

# Copy binary from backend builder
COPY --from=backend-builder /app/main .

# Copy frontend build from frontend builder
COPY --from=frontend-builder /app/frontend/dist ./static

# Copy migrations
COPY --from=backend-builder /app/internal/infra/db/migrations ./migrations/

# Set ownership
RUN chown -R appuser:appuser /app

USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
