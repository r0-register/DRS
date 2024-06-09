# 使用带有 Alpine 的 Go 官方镜像作为构建阶段基础镜像
FROM golang:alpine as builder

# 启用 Go modules，配置国内代理以加快依赖下载速度
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /workspace

# 复制 Go 模块的描述文件
COPY go.mod go.sum ./

# 下载依赖项
RUN go mod download

# 复制源代码文件，特别注意目录路径
COPY pkg/ ./pkg/

# 构建应用程序，确保指定正确的路径到 main.go
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o drs-scheduler ./pkg/main.go

# 使用轻量级的 Alpine 镜像作为最终镜像
FROM alpine:latest

# 安装时区数据和CA证书
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的应用
COPY --from=builder /workspace/drs-scheduler .

# 设置执行权限
RUN chmod +x drs-scheduler

# 指定容器启动时运行的命令
CMD ["./drs-scheduler"]
