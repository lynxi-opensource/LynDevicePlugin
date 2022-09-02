# LynDevicePlugin

## 使用

```sh
# 安装LynContainer和docker
# push镜像
make push-lynxi-k8s-device-plugin
make push-lynxi-exporter
cd apu-feature-discovery
make push
# 构建chart
make chart
```

## 更新版本

1. 更新build.Dockerfile中的基础sdk镜像
2. 更新Makefile中的版本号
3. 更新LynDevicePlugin中的image tag和版本号
4. 构建发布
