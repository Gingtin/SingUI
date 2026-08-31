# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app/backend
RUN apk add --no-cache git gcc musl-dev
COPY backend/go.mod backend/go.sum* ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend-builder /app/backend/cmd/server/dist ./cmd/server/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o singbox-ui ./cmd/server

# Stage 3: Final runtime
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata curl \
    && mkdir -p /app/data /app/config

# Download sing-box core
ARG TARGETARCH
RUN case "${TARGETARCH}" in \
        "amd64") SB_ARCH="linux-amd64" ;; \
        "arm64") SB_ARCH="linux-arm64" ;; \
        *) SB_ARCH="linux-amd64" ;; \
    esac \
    && curl -fsSL "https://github.com/SagerNet/sing-box/releases/download/v1.9.7/sing-box-1.9.7-${SB_ARCH}.tar.gz" -o /tmp/sb.tar.gz \
    && tar -xzf /tmp/sb.tar.gz -C /tmp/ \
    && cp /tmp/sing-box-*/sing-box /usr/local/bin/sing-box \
    && chmod +x /usr/local/bin/sing-box \
    && rm -rf /tmp/*

COPY --from=backend-builder /app/backend/singbox-ui /app/singbox-ui
RUN chmod +x /app/singbox-ui

EXPOSE 2096 9090
VOLUME ["/app/data", "/app/config"]

ENTRYPOINT ["/app/singbox-ui", "-d", "/app/data/singbox_ui.db", "-p", "2096"]
