# 集群模式下微服务Prometheus监控

Prometheus单机部署，新增了抓取配置后，需要重启premetheus服务，集群模式下推荐使用PrometheusOperator。

安装Promethues instance
```
kubectl apply -f instance.yaml
```

创建服务监听
```
kubectl apply -f monitors
```

