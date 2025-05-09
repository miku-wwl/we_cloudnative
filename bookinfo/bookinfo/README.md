
## 初始化项目结构

[cobra](https://github.com/spf13/cobra)安装
```
go get -u github.com/spf13/cobra@v1.8.0
go install github.com/spf13/cobra-cli@latest
# 配置gobin到环境变量 $GOBIN根据自己的go env gobin对应的地址去做替换
export PATH=$PATH:$GOBIN
```

## 初始化web模块 Gin

安装[Gin](https://github.com/gin-gonic/gin)
```
go get -u github.com/gin-gonic/gin@v1.10.0
```


安装[Gin](https://github.com/gin-gonic/gin)
```
go get -u github.com/gin-gonic/gin@v1.10.0
```

安装[Gin指标采集](github.com/zsais/go-gin-prometheus)

`$ go get github.com/zsais/go-gin-prometheus`

`demo.go`
```
package main

import (
  "github.com/gin-gonic/gin"
  "github.com/zsais/go-gin-prometheus"
)

func main() {
  r := gin.New()

  p := ginprometheus.NewPrometheus("gin")
  p.Use(r)

  r.GET("/", func(c *gin.Context) {
    c.JSON(200, "Hello world!")
  })

  r.Run(":29090")
}
```



## 功能开发
- [x] 首页（展示服务间的调用关系）
- [x] 图书详情页


## Bookinfo docker构建部署
在bookinfo目录下使用指令
docker build -t bookinfo:v0.1 .
docker run -p8002:8080 -d  bookinfo:v0.1
curl localhost:8002

https://goharbor.io/docs/2.12.0/install-config/
gencert.sh
```
openssl genrsa -out ca.key 4096

openssl req -x509 -new -nodes -sha512 -days 3650 \
 -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=MyPersonal Root CA" \
 -key ca.key \
 -out ca.crt

openssl genrsa -out harbor.com.key 4096

openssl req -sha512 -new \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=harbor.com" \
    -key harbor.com.key \
    -out harbor.com.csr


cat > v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=harbor.com
DNS.2=yourdomain
DNS.3=hostname
EOF

openssl x509 -req -sha512 -days 3650 \
    -extfile v3.ext \
    -CA ca.crt -CAkey ca.key -CAcreateserial \
    -in harbor.com.csr \
    -out harbor.com.crt


mkdir -p /data/cert/
cp harbor.com.crt /data/cert/
cp harbor.com.key /data/cert/

openssl x509 -inform PEM -in harbor.com.crt -out harbor.com.cert


cp harbor.com.cert /etc/docker/certs.d/harbor.com/
cp harbor.com.key /etc/docker/certs.d/harbor.com/
cp ca.crt /etc/docker/certs.d/harbor.com/

systemctl restart docker
```


sudo yum install tree
tree
└── harbor.com
    ├── ca.crt
    ├── harbor.com.cert
    └── harbor.com.key




cp harbor.yml harbor.yml.tmpl
docker load < harbor.v2.13.0.tar.gz

vim harbor.yml

ls /data/cert/

修改harbor.yml 
```
[root@192 harbor]# vim harbor.yml
# Configuration file of Harbor

# The IP address or hostname to access admin UI and registry service.
# DO NOT use localhost or 127.0.0.1, because Harbor needs to be accessed by external clients.
hostname: harbor.com

# http related config
http:
  # port for http, default is 80. If https enabled, this port will redirect to https port
  port: 80

# https related config
https:
  # https port for harbor, default is 443
  port: 443
  # The path of cert and key files for nginx
  certificate: /data/cert/harbor.com.crt
  private_key: /data/cert/harbor.com.key
```

docker compose up -d


vim harbor.yml 查看初始密码
![alt text](image.png)

## Harbor镜像推送与客户端证书配置
vim /etc/hosts
末尾添加 192.168.171.133 harbor.com


配置文件永久跳过
```
# 1. 编辑 Docker 配置文件
sudo nano /etc/docker/daemon.json

# 2. 添加私有仓库到不安全列表（多行用逗号分隔）
{
  "insecure-registries": ["harbor.com"],
}

# 3. 重启 Docker 服务
sudo systemctl restart docker
```

docker login harbor.com
admin
Harbor12345

docker push  harbor.com/test/bookinfo:v0.1


## Docker 容器与镜像的相互转换
docker commit id bookinfo:v0.2




## Harbor 服务运行状态监控
cd /opt/harbor/harbor/
vim harbor.yml
```
metric:
  enabled: true
  port: 9095
  path: /metrics
```
./prepare
vim docker-compose.yml 查看9095端口是否启动
docker compose up -d
curl localhost:9095/metrics


## Harbor 访问失败定位
通过docker ps | grep harbor 发现 harbor 启动的container很少
使用docker logs -f container_id

2025-04-14T16:06:08Z [ERROR] [/lib/cache/cache.go:126]: failed to ping redis://redis:6379?idle_timeout_seconds=30, retry after 1.734749814s : dial tcp 104.18.25.196:6379: i/o timeout

发现以来的redis未启动

解决方案
进入harbor 目录，重启harbor即可


## 理解声明式API特性

## 理解声明式API特性

## kubernetes安装选型

## MiniKube环境搭建

## kubernetes核心技能提炼


