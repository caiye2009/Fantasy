FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# 1. 安装构建依赖
RUN apk add --no-cache git build-base libc6-compat

# 2. 先复制 go.mod 和 go.sum，下载依赖（可以缓存）
COPY go.mod go.sum ./
RUN go mod download

# 3. 复制剩余代码
COPY . .

# 4. 编译
RUN make build

# =======================
# 最小运行镜像
FROM alpine:3.18
WORKDIR /app

# 复制构建好的二进制
COPY --from=builder /app/fantasy-back .

EXPOSE 8080
CMD ["./fantasy-back"]
