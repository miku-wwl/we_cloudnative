# buildx 使用指南

以下操作在linux/amd64环境演示。

## 开启docker实验特性

```
vim /etc/docker/daemon.json
```
添加配置：
```json
{
    ...
    "experimental": true,
    ...
}
```
重启docker
```
systemct restart docker
```
如果重启之后发现 `docker buildx` 命令还是没有，则可以考虑直接下载 [buildx](https://github.com/docker/buildx)

## 跨架构配置
例如，在 x86_64 主机上构建和测试适用于 ARM 架构的镜像。

安装qemu
```bash
sudo apt install -y qemu-user-static binfmt-support
```

启用qemu
```bash
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

## 创建buildx构建实例

默认已经有了一个构建实例，但是默认的构建器不支持并发构建镜像。
```
# 创建
docker buildx create --driver-opt network=host --platform=linux/arm64,linux/amd64 --name bk-builder
# 使用
docker buildx use bk-builder
# 初始化
docker buildx inspect --bootstrap
```

设置ca证书
> 课程演示的是自签名证书，因此需要将自签名证书拷贝到构建器中。
```
builder=$(docker ps | grep buildx | awk '{print $1}')
docker cp /etc/docker/certs.d/harbor.imooc.com/ca.crt $builder:/usr/local/share/ca-certificates/ 
docker exec $builder update-ca-certificates
docker restart $builder 
```

其他命令：
```
# 查询构建器
docker buildx ls
# 删除构建器
docker buildx rm bk-builder
```