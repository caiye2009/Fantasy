# 使用 Alpine 版 Node 镜像
FROM node:22-alpine AS builder

# 设置环境变量
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV NODE_OPTIONS=--max-old-space-size=8192
ENV TZ=Asia/Shanghai
ENV CI=true

# 安装 corepack 和其他依赖
RUN npm i -g corepack

# 安装解压工具（Alpine 使用 apk 安装）
RUN apk add --no-cache unzip

WORKDIR /app

# 复制项目文件
COPY . /app

# 安装依赖
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

# build web-ele
RUN pnpm --filter=apps/web-ele run build

# 解压 dist.zip
RUN unzip -o apps/web-ele/dist.zip -d apps/web-ele/dist

RUN echo "Builder Success 🎉"

# 使用 Alpine 版的 nginx 镜像作为生产环境镜像
FROM nginx:stable-alpine AS production

# 配置 nginx
RUN echo "types { application/javascript js mjs; }" > /etc/nginx/conf.d/mjs.conf \
    && rm -rf /etc/nginx/conf.d/default.conf

# 复制构建产物
COPY --from=builder /app/apps/web-ele/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY --from=builder /app/scripts/deploy/nginx.conf /etc/nginx/nginx.conf

EXPOSE 8080

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
