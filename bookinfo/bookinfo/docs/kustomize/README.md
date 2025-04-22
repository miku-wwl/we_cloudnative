# kustomize 使用手册
Kustomize 是一种配置管理解决方案，在不改变原始yaml文件的情况下，通过patch有选择地覆盖来生成新的配置，很方便维护多版本的yaml文件，特别是项目较大或者需要处理很多case情况下，非常方便并且易于维护。

k8s资源生成器

## 生成资源

生成一个  application.properties 文件
```
cat <<EOF >application.properties
FOO=Bar
EOF
```

将 application.properties 作为key
```
cat <<EOF >./kustomization.yaml
configMapGenerator:
- name: example-configmap-1
  files:
  - application.properties
EOF
```
根据键值对展示

```
cat <<EOF >./kustomization.yaml
configMapGenerator:
- name: example-configmap-2
  literals:
  - FOO=Bar
EOF
```

## 演示基准与overlay

演示secret动态创建
资源批量替换
组织，定时，补丁
