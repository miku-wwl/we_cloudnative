# istio安装

istio安装需要4c8G的基础环境，在服务网格特性之前的学习，不需要这么高的配置，所以采用的是默认配置（2c4G）。

## 删除集群并重新创建

```
minikube delete
# mac环境，windows环境推荐virtualbox
minikube start --cpus=4 --memory=8192 --driver=qemu2 --network socket_vmnet
```

## 安装istio

安装istio operator
```
minikube addons enable istio-provisioner
```

安装istiod和istio gateway
```
minikube addons enable istio
```

检查istio-system和istio-operator的pod都处于running状态表示安装成功。

## kubeadm集群安装istio

下载目前最新版安装文件
https://github.com/istio/istio/releases/tag/1.25.2

tar zxvf istio-1.25.2-linux-amd64.tar.gz

cd istio-1.25.2

export PATH=$PWD/bin:$PATH

istioctl install -y

如果没有负载均衡器, loadbalancer 需要修改istio-ingressgateway 的type为NodePort

kubectl edit svc istio-ingressgateway -n istio-system
进入修改即可


卸载istio

istioctl uninstall -y --purge