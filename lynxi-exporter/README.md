# lynxi-exporter

lynxi-exporter用于演示如何实现提供Prometheus格式的apu指标数据，并使用DaemonSet和Prometheus-Operator自动化部署。

## 实现方法

1. [main.go](main.go), [Dockerfile](Dockerfile), [build_push.sh](build_push.sh): 基于lynxi-image环境，调用lynxi-smi，并使用Prometheus提供的包定义metrics http接口，提供apu_total_usage和apu_count指标。
2. [DaemonSet.yml](DaemonSet.yml), [Service.yml](Service.yml): 为带有`lynxi.com=apu`的节点部署DaemonSet，并部署Service使得其它应用如Prometheus能够访问。
3. [ServiceMonitor.yml](ServiceMonitor.yml): 创建ServiceMonitor资源，用于Prometheus-operator自动配置Prometheus，Prometheus监听到配置文件改变后重新加载配置文件，并提供apu指标的查询接口。
4. apu: 使用kubesphere提供的自定义监控界面创建折线图表。

## 部署方法

```sh
make service_monitor
```

## 指标示例

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
