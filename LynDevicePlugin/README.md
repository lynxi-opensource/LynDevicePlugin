# LynDevicePlugin

一个用于安装apu-feature-discovery, lynxi-k8s-device-plugin, lynxi-exporter的helm chart.

- `apu-feature-discovery`：自动发现设备
- `lynxi-k8s-device-plugin`：用于对接lynxi设备和k8s集群
- `lynxi-exporter`：提供Prometheus格式的数据
- 可选的Prometheus-operator自定义资源
  - `lynxi-expoter-service`：lynxi-expoter的service
  - `lynxi-expoter-service-monitor`：用于Prometheus-operator的ServiceMonitor

## 前置条件

1. k8s集群

## 配置插有lynxi板卡的k8s节点

1. 安装LynDriver请参考sdk用户手册
2. 安装lynxi-docker请参考lynxi-docker用户手册
3. 增加改节点到k8s集群

## 安装LynDevicePlugin

1. 创建namespace：`kubectl create namespace project-system`
2. 安装helm，请[参考](https://helm.sh/docs/intro/quickstart/)
3. 安装LynDevicePlugin:
   1. 如果k8s集权中没有安装Prometheus-operator，执行`helm install -n project-system --set lynxiExporterServiceMonitor.enable=false lynxi-device-plugin LynDevicePlugin-0.1.0.tgz`安装
   2. 否则，执行`helm install -n project-system lynxi-device-plugin LynDevicePlugin-0.1.0.tgz`安装
4. 查看安装是否成功: 执行`helm list -n project-system`，查看lynxi-device-plugin的状态是否为deployed

```
NAME                    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                           APP VERSION
LynDevicePlugin      project-system  1               2022-01-26 10:04:28.7551467 +0800 CST   deployed        LynDevicePlugin-1.0.0        1.16.0
```

## 卸载

1. `helm uninstall LynDevicePlugin`

## lynxi-exporter提供的Prometheus指标

lynxi-exporter的/metircs输出示例：

```yaml
# HELP lynxi_board_power The power of the board with unit mW
# TYPE lynxi_board_power gauge
lynxi_board_power{BoardID="1",Manufacturer="lynxi.com",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 14358
lynxi_board_power{BoardID="1",Manufacturer="lynxi.com",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 13211
# HELP lynxi_device_apu_usage The apu usage of the device with unit %
# TYPE lynxi_device_apu_usage gauge
lynxi_device_apu_usage{BoardID="1",ID="0",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_apu_usage{BoardID="1",ID="1",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_apu_usage{BoardID="1",ID="2",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_apu_usage{BoardID="1",ID="3",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_apu_usage{BoardID="1",ID="4",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_apu_usage{BoardID="1",ID="5",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
# HELP lynxi_device_arm_usage The arm usage of the device with unit %
# TYPE lynxi_device_arm_usage gauge
lynxi_device_arm_usage{BoardID="1",ID="0",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 10
lynxi_device_arm_usage{BoardID="1",ID="1",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 19
lynxi_device_arm_usage{BoardID="1",ID="2",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 20
lynxi_device_arm_usage{BoardID="1",ID="3",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_arm_usage{BoardID="1",ID="4",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_arm_usage{BoardID="1",ID="5",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 28
# HELP lynxi_device_current_temp The current temperature of the device with unit ℃
# TYPE lynxi_device_current_temp gauge
lynxi_device_current_temp{BoardID="1",ID="0",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 63
lynxi_device_current_temp{BoardID="1",ID="1",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 62
lynxi_device_current_temp{BoardID="1",ID="2",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 63
lynxi_device_current_temp{BoardID="1",ID="3",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 54
lynxi_device_current_temp{BoardID="1",ID="4",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 54
lynxi_device_current_temp{BoardID="1",ID="5",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 54
# HELP lynxi_device_ipe_usage The ipe usage of the device with unit FPS
# TYPE lynxi_device_ipe_usage gauge
lynxi_device_ipe_usage{BoardID="1",ID="0",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_ipe_usage{BoardID="1",ID="1",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_ipe_usage{BoardID="1",ID="2",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_ipe_usage{BoardID="1",ID="3",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_ipe_usage{BoardID="1",ID="4",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_ipe_usage{BoardID="1",ID="5",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
# HELP lynxi_device_mem_used The memory used of the device with unit KB
# TYPE lynxi_device_mem_used gauge
lynxi_device_mem_used{BoardID="1",ID="0",Manufacturer="lynxi.com",MemTotal="8388608",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 4.472832e+06
lynxi_device_mem_used{BoardID="1",ID="1",Manufacturer="lynxi.com",MemTotal="8388608",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 4.472832e+06
lynxi_device_mem_used{BoardID="1",ID="2",Manufacturer="lynxi.com",MemTotal="8388608",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 4.472832e+06
lynxi_device_mem_used{BoardID="1",ID="3",Manufacturer="lynxi.com",MemTotal="8388608",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 4.472832e+06
lynxi_device_mem_used{BoardID="1",ID="4",Manufacturer="lynxi.com",MemTotal="8388608",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 4.472832e+06
lynxi_device_mem_used{BoardID="1",ID="5",Manufacturer="lynxi.com",MemTotal="8388608",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 4.472832e+06
# HELP lynxi_device_state The state of the device
# TYPE lynxi_device_state gauge
lynxi_device_state{BoardID="1",ID="0",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 1
lynxi_device_state{BoardID="1",ID="1",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 1
lynxi_device_state{BoardID="1",ID="2",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 1
lynxi_device_state{BoardID="1",ID="3",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 1
lynxi_device_state{BoardID="1",ID="4",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 1
lynxi_device_state{BoardID="1",ID="5",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 1
# HELP lynxi_device_vic_usage The vic usage of the device with unit %
# TYPE lynxi_device_vic_usage gauge
lynxi_device_vic_usage{BoardID="1",ID="0",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_vic_usage{BoardID="1",ID="1",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_vic_usage{BoardID="1",ID="2",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00045"} 0
lynxi_device_vic_usage{BoardID="1",ID="3",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_vic_usage{BoardID="1",ID="4",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
lynxi_device_vic_usage{BoardID="1",ID="5",Manufacturer="lynxi.com",Model="KA200",MountTime="2022-05-27T08:40:19Z",ProductName="HP300",SerialNumber="TICLX21B23A00048"} 0
# HELP lynxi_exporter_state Is there any error in lynxi-exporter internal, please see the logs, will auto recovering after a while. 1 is ok, 0 is err.
# TYPE lynxi_exporter_state gauge
lynxi_exporter_state 1
# HELP lynxi_pod_container_device_count The device ids and number of devices for each pod container.
# TYPE lynxi_pod_container_device_count gauge
lynxi_pod_container_device_count{device_ids="1",owner_container="nlp5962",owner_namespace="dnn-project",owner_pod="nlp5962-66589469b9-mzbfc"} 1
lynxi_pod_container_device_count{device_ids="2",owner_container="ar1125",owner_namespace="dnn-project",owner_pod="ar1125-5575c9f444-w9qgs"} 1
lynxi_pod_container_device_count{device_ids="3",owner_container="ir5613",owner_namespace="dnn-project",owner_pod="ir5613-59dd95fccf-5b8jf"} 1
lynxi_pod_container_device_count{device_ids="5",owner_container="ir1877",owner_namespace="dnn-project",owner_pod="ir1877-747574c6b7-8rprb"} 1
```
