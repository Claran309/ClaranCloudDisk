# 构建阶段
FROM golang:1.25.3-alpine AS builder

# 安装必要的构建工具
RUN apk --no-cache add git

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（静态链接，减小依赖）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./main.go

# 生成 Swagger 文档
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# 运行阶段 - 使用更小的基础镜像
FROM alpine:3.19

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata && \
    mkdir -p /app && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=builder --chown=appuser:appgroup /app/docs ./docs
COPY --from=builder --chown=appuser:appgroup /app/config ./config

# 创建必要的目录并复制默认头像
RUN mkdir -p ./Avatars/DefaultAvatar
COPY --from=builder --chown=appuser:appgroup /app/Avatars/DefaultAvatar/DefaultAvatar.png ./Avatars/DefaultAvatar/

# 确保配置文件存在
RUN test -f ./config/config.yaml || (echo "Config file missing!" && exit 1)

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 运行应用
ENTRYPOINT ["./main"]