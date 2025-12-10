FROM node:22-slim AS builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV NODE_OPTIONS=--max-old-space-size=8192
ENV TZ=Asia/Shanghai
ENV NODE_ENV=production

RUN npm i -g corepack

WORKDIR /app

# 复制所有文件
COPY . /app

ENV CI=true

# 安装依赖
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

# 构建 web-ele 项目（生产环境）
RUN pnpm run build:ele

# 验证构建产物
RUN ls -lah /app/apps/web-ele/dist
RUN echo "Builder Success 🎉"

FROM nginx:stable-alpine AS production

# 配置 nginx
RUN echo "types { application/javascript js mjs; }" > /etc/nginx/conf.d/mjs.conf \
    && rm -rf /etc/nginx/conf.d/default.conf

# 复制构建产物
COPY --from=builder /app/apps/web-ele/dist /usr/share/nginx/html

# 验证文件
RUN ls -la /usr/share/nginx/html

# 复制 nginx 配置
COPY --from=builder /app/scripts/deploy/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]