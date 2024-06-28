package podresources

import (
	"context"
	"lyndeviceplugin/utils/singleflight"
	"net"
	"time"

	"google.golang.org/grpc"
	podresources "k8s.io/kubelet/pkg/apis/podresources/v1"
)

type retType struct {
	ret []Resource
	err error
}

type PodResources struct {
	conn   *grpc.ClientConn
	client podresources.PodResourcesListerClient
	sf     singleflight.Singleflight[retType]
}

func New() (*PodResources, error) {
	conn, err := dial("/var/lib/kubelet/pod-resources/kubelet.sock", 3*time.Second)
	if err != nil {
		return nil, err
	}

	client := podresources.NewPodResourcesListerClient(conn)
	return &PodResources{conn, client, singleflight.New[retType]()}, nil
}

func (m *PodResources) Close() error {
	return m.conn.Close()
}

type ResourceOwner struct {
	Pod       string
	Namespace string
	Container string
}

type Resource struct {
	ResourceOwner
	IDs []string
}

func (m *PodResources) get() (ret []Resource, err error) {
	resp, err := m.client.List(context.Background(), &podresources.ListPodResourcesRequest{})
	if err != nil {
		return
	}
	for _, pod := range resp.GetPodResources() {
		for _, container := range pod.GetContainers() {
			res := Resource{
				ResourceOwner{pod.GetName(), pod.GetNamespace(), container.GetName()}, nil,
			}
			for _, device := range container.GetDevices() {
				if device.GetResourceName() == "lynxi.com/device" {
					res.IDs = append(res.IDs, device.GetDeviceIds()...)
				}
			}
			if len(res.IDs) > 0 {
				ret = append(ret, res)
			}
		}
	}
	return
}

func (m *PodResources) Get() ([]Resource, error) {
	ret := m.sf.Fly(func() retType {
		ret, err := m.get()
		return retType{ret, err}
	})
	return ret.ret, ret.err
}

// dial establishes the gRPC communication with the registered device plugin.
func dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
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
