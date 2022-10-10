package metrics

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	podresources "k8s.io/kubelet/pkg/apis/podresources/v1"
)

var _ Recorder = &PodContainerRecorder{}

// States 定义和记录所有状态相关的Prometheus指标
type PodContainerRecorder struct {
	lynxiPodContainerDeviceCount *prometheus.GaugeVec
	timeout                      time.Duration
}

// NewStatesRecorder 构造一个StatesRecorder并初始化指标
func NewPodContainerRecorder(timeout time.Duration) *PodContainerRecorder {
	ret := &PodContainerRecorder{
		lynxiPodContainerDeviceCount: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_pod_container_device_count",
			Help: "The device ids and number of devices for each pod container.",
		}, labelsForPodContainer()),
		timeout: timeout,
	}
	return ret
}

func labelsForPodContainer() []string {
	return []string{"owner_pod", "owner_container", "owner_namespace", "device_ids"}
}

// Record 一直阻塞不会返回错误，外部通过lynxi_exporter_state或日志查看exporter的状态是否正常
func (m *PodContainerRecorder) Record() error {
	for {
		log.Println("connect kubelet")
		conn, err := m.dial("/var/lib/kubelet/pod-resources/kubelet.sock", 5*time.Second)
		if err != nil {
			GlobalRecorder.logError(err)
			continue
		}
		defer conn.Close()

		client := podresources.NewPodResourcesListerClient(conn)
		ticker := time.NewTicker(m.timeout)
		for range ticker.C {
			resp, err := client.List(context.Background(), &podresources.ListPodResourcesRequest{})
			if err != nil {
				GlobalRecorder.logError(err)
				break
			} else {
				m.lynxiPodContainerDeviceCount.Reset()
				for _, pod := range resp.PodResources {
					for _, container := range pod.Containers {
						for _, device := range container.Devices {
							if device.ResourceName == "lynxi.com/device" {
								m.lynxiPodContainerDeviceCount.WithLabelValues(
									pod.Name, container.Name, pod.Namespace,
									strings.Join(device.DeviceIds, ",")).Set(float64(len(device.DeviceIds)))
							}
						}
					}
				}
			}
		}
	}
}

// dial establishes the gRPC communication with the registered device plugin.
func (m *PodContainerRecorder) dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}
