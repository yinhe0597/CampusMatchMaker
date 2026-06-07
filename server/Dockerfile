# ===== 构建阶段 =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安装 git（某些依赖需要）
RUN apk --no-cache add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server cmd/main.go

# ===== 运行阶段 =====
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata curl
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /server .
COPY configs/ ./configs/
COPY migrations/ ./migrations/

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD curl -sf http://localhost:8080/api/v1/health || exit 1

CMD ["./server"]
