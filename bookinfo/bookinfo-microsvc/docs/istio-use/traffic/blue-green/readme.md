# 蓝绿发布

蓝：代表稳定，旧的系统
绿：代表新的，新的系统

需要两套系统，流量随意切换。

## 怎么实践

- 修改代码 reviews，reviews 显示podname
- 部署服务的时候，不同版本通过名称表示，例如reviews-v1 ,reveiws-v2