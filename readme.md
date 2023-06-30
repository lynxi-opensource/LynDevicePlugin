# LynDevicePlugin

## 使用

```sh
# 安装docker、golang1.18+、rust1.60+
# 构建和推送镜像
./build.sh

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

1. 更新Makefile和build.sh中的版本号
2. 更新LynDevicePlugin中的image tag和版本号
3. 构建发布
