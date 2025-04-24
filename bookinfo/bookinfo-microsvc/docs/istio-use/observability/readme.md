# 服务网格可观测性

整合前面学习的prometheus-operator，jaeger，grafana，新增loki和kiali。

可观测性三个指标：metrics、trace、log.

## prometheus+grafana监控

安装prometheus oprator，安装grafana。
```
kubectl create -f promoperator/bundle.yaml
```
创建实例
```
kubectl apply -f promoperator/instance.yaml
```
创建采集目标
```
kubectl apply -f promoperator/monitor.yaml
```
安装grafana监控
```
kubectl apply -f grafana/grafana.yaml
```

## jaeger调用链路追踪
重新配置istio，注入jaeger配置
istio基于istio-operator安装的，可以通过命令查询：
```
kubectl get istiooperator -A
```

默认的istio没有任何配置，执行如下命令，指定jaeger配置：
```
kubectl apply -f jaeger/istio-operator.yaml
```
重新启动注入了sidecar的服务。

## loki日志监控
安装日志收集服务端
```
kubectl apply -f loki/loki.yaml
```
安装日志收集客户端
```
kubectl apply -f loki/promtai.yaml
```
grafana查询日志
## kiali服务网格观测
kiali用于监控和管理 Istio 服务网格。它提供了一个直观的界面来查看和分析微服务之间的流量、健康状况和性能。
安装
```
kubectl apply -f kiali/kaili.yaml
```