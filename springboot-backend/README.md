# springboot-backend

## Dockerfile构建springboot镜像，本地无须安装maven

```shell
使用构建指令并执行
docker build -t springboot-backend:v0.1 .
docker run -itd -p8081:8080 springboot-backend:v0.1
```