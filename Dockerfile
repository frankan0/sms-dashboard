# Stage 1: Build the Go backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

# 设置国内 Go 代理加速下载
ENV GOPROXY=https://goproxy.cn,direct

# Install git if needed for dependencies
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# The local 'web/dist' directory must already exist (built by 'npm run build' locally)
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o sms-dashboard ./cmd/server/main.go

# Stage 3: Final production image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /app
# Create data directory for SQLite persistence
RUN mkdir -p /app/data
COPY --from=backend-builder /app/sms-dashboard .

# Set default environment variables
ENV PORT=8080
ENV API_TOKEN=default-api-token
ENV JWT_SECRET=default-jwt-secret
ENV DB_PATH=/app/data/sms.db
ENV GIN_MODE=release
ENV SECRET=defaultSecret

# Expose port
EXPOSE 8080

# Volume for persistence
VOLUME /app/data

# Run the app
CMD ["./sms-dashboard"]
