# jaeger-operator安装

kubernetes集群模式下使用jaeger，安装jaeger-operator

[安装cert-manger](https://github.com/cert-manager/cert-manager)
```
kubectl apply -f cert-manager.yaml
```
> 在 Kubernetes 集群中自动申请、颁发、续订和管理 TLS/SSL 证书。

安装证书后确认
```
kubectl get pod -n cert-manager
```
输出如下内容表示安装成功了
```
NAME                                       READY   STATUS    RESTARTS       AGE
cert-manager-6796d554c5-vq7zc              1/1     Running   26 (66m ago)   43h
cert-manager-cainjector-77cd756b5d-7h9nz   1/1     Running   2 (26h ago)    43h
cert-manager-webhook-dbb5879d7-9k28s       1/1     Running   2 (26h ago)    43h
```

[安装Jaeger](https://www.jaegertracing.io/docs/1.62/operator/#tracing-and-debugging-the-operator)

```
kubectl create namespace observability
kubectl create -f jaeger-operator.yaml
```

创建实例
```
kubectl apply -f sample.yaml
```
默认会创建service，ingress等资源，通过ingress即可访问。
