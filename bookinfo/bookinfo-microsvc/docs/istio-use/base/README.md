# istio使用


## sidecar注入

```
kubectl label namespace test istio-injection=enabled
```

## 安装gateway
查看istio自定义资源支持的版本
查看gateway
```
kubectl get crd gateways.networking.istio.io -o json | jq -r '.spec.versions[].name'
```
查看virtualservice
```
kubectl get crd virtualservices.networking.istio.io -o json | jq -r '.spec.versions[].name'
```

