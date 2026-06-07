# ===== 构建阶段 =====
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/main.go

# ===== 运行阶段 =====
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /server .
COPY configs/ ./configs/
COPY migrations/ ./migrations/

EXPOSE 8080

CMD ["./server"]
