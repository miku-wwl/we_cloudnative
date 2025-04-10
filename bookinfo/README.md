
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

## 功能开发
- [x] 首页（展示服务间的调用关系）
- [x] 图书详情页