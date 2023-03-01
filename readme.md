# LynDevicePlugin

## 使用

```sh
# 安装LynContainer和docker
# 构建amd64镜像，在hp300、he200环境
make build-amd64 push-amd64
cd apu-feature-discovery
make push build-amd64 push-amd64
cd ..

# 构建arm64镜像，在hs110上执行
make build-arm64 push-arm64
cd apu-feature-discovery
make push build-arm64 push-arm64
cd ..

# 创建docker manifest
make docker-manifest
cd apu-feature-discovery
make docker-manifest
cd ..

# 构建chart
make chart
# 安装
make install
# 查看
make list
# 卸载
make uninstall
```

## 更新版本

1. 更新Makefile中的版本号
2. 更新LynDevicePlugin中的image tag和版本号
3. 构建发布
