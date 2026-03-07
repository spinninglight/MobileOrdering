# 使用最新的 Go 1.25 镜像
FROM golang:1.25-alpine AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app

# 先拷贝依赖文件，利用 Docker 缓存层优化构建速度
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 拷贝源代码
COPY backend/ .

# 关键点：指定 main.go 所在的子目录进行编译
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

# 运行阶段
FROM alpine:3.21
WORKDIR /root/

# 拷贝编译后的二进制文件
COPY --from=builder /app/server .

# 如果你有 .env 文件或静态资源，也需要拷贝
# COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./server"]