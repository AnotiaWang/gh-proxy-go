# 使用官方 Go 基础镜像，这里可以指定 Go 版本
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

RUN apk --no-cache add ca-certificates

# 复制 go.mod 和 go.sum 文件
COPY . .
# 下载依赖项
RUN go mod download

# 构建可执行文件
# CGO_ENABLED=0 用于静态链接，以便在 scratch 镜像中运行
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用 scratch 作为基础镜像以减小镜像大小
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 从 builder 镜像中复制构建好的可执行文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 80

# 运行可执行文件
CMD ["./main"]
