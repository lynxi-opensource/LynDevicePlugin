# LynDevicePlugin

## 使用

```sh
# 安装LynContainer和docker
# push镜像
make push
cd apu-feature-discovery
make push
cd ..

# push hs110-kubeedge镜像，需要在hs110上执行
make hs110-kubeedge-push

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
