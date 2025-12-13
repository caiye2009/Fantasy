# ---- 依赖缓存阶段 ----
FROM golang:1.25.4-alpine AS base

# 国内 apk 源加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /app

# Go 代理（你现在的网络情况必须保留）
ENV GOPROXY=https://goproxy.cn,direct

# 只复制 go.mod 和 go.sum，用于依赖缓存
COPY go.mod go.sum ./

# 基础依赖 + Go modules
RUN apk update && apk add --no-cache \
    git \
    build-base \
    libc6-compat \
    tzdata \
    && go mod download

# ---- 编译阶段 ----
FROM base AS builder

WORKDIR /app

# 复制完整项目（包含 config/）
COPY . .

# 使用 Makefile 编译（你明确说了：用 make）
RUN --mount=type=cache,target=/root/.cache/go-build make

# ---- 最终运行镜像 ----
FROM alpine:3.18

WORKDIR /app

# 运行依赖
RUN apk add --no-cache libc6-compat tzdata

# 时区（解决 Asia/Shanghai）
ENV TZ=Asia/Shanghai
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# ⭐ 复制二进制
COPY --from=builder /app/fantasy-back ./

# ⭐⭐ 关键：复制 config 目录（与你现在的代码路径完全一致）
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./fantasy-back"]