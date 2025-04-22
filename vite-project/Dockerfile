FROM node:22.14.0-alpine

WORKDIR /app

COPY package.json package.json
COPY package-lock.json package-lock.json

# 设置 npm 镜像源
RUN npm config set registry https://registry.npmmirror.com

RUN npm install

COPY . /app

ENTRYPOINT ["npm", "run", "dev"]

EXPOSE 5173


# 修改vite.config.ts 文件

# export default defineConfig({
#   plugins: [react()],
#   server: {
#     host: '0.0.0.0'
#   }
# })

