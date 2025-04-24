# 超时&重试

超时：我们流量管理的方式动态超时时间。
重试：例如丢包、偶发报错等情况，发起重试。

## 如何实践

- 模拟ratings timeout，修改reviews的 timeout=1,reviews的default timeout将会被重置。
- 模拟ratings timeout，，修改reviews的 timeout=1 reviews 5xx状态码,reviews的表现形式超时啦，productpage发起reviews请求，验证超时。