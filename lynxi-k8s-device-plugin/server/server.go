// package server 负责与kubelet建立通信
package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func newOSWatcher(sigs ...os.Signal) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sigs...)

	return sigChan
}

// Server 描述了与kubelet建立通信所需的接口和参数
type Server interface {
	// Serve 启动server，在server启动后注册到kubelet，并阻塞到server退出
	Run(socket, resourceName string, service pluginapi.DevicePluginServer) error
}

type ServerImp struct {
	Crash <-chan error
}

//
func (m *ServerImp) Run(socket, resourceName string, service pluginapi.DevicePluginServer) error {
	log.Println("Starting OS watcher.")
	sigs := newOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var srv *LynxiServer

restart:
	if srv != nil {
		srv.Stop()
	} else {
		srv = &LynxiServer{
			resourceName: resourceName,
			socket:       path.Join(pluginapi.DevicePluginPath, socket),
			// These will be reinitialized every
			// time the plugin server is restarted.
			server: nil,
			plugin: &service,
		}
	}

	if err := srv.Start(); err != nil {
		log.Println(err)
		goto restart
	}

events:
	for {
		select {
		// If there was an error starting any plugins, restart them all.
		case <-m.Crash:
			goto restart
		case s := <-sigs:
			switch s {
			case syscall.SIGHUP:
				log.Println("Received SIGHUP, restarting.")
				goto restart
			default:
				log.Printf("Received signal \"%v\", shutting down.", s)
				srv.Stop()
				break events
			}
		}
	}

	return nil
}

type LynxiServer struct {
	socket       string                        // local grpc socket file name
	resourceName string                        // 用于k8s分配资源时使用的资源名
	server       *grpc.Server                  // grpc server
	plugin       *pluginapi.DevicePluginServer // k8s plugin
}

//
func (m *LynxiServer) Start() error {
	m.initialize()

	err := m.Serve()
	if err != nil {
		log.Printf("Could not start device plugin for '%s': %s", m.resourceName, err)
		m.cleanup()
		return err
	}
	log.Printf("Starting to serve '%s' on %s", m.resourceName, m.socket)

	err = m.Register()
	if err != nil {
		log.Printf("Could not register device plugin: %s", err)
		m.Stop()
		return err
	}
	log.Printf("Registered device plugin for '%s' with Kubelet", m.resourceName)
	return nil
}

func (m *LynxiServer) Stop() error {
	if m == nil || m.server == nil {
		return nil
	}
	log.Printf("Stopping to serve '%s' on %s", m.resourceName, m.socket)
	m.server.Stop()

	err := os.Remove(path.Join(pluginapi.DevicePluginPath, m.socket))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	m.cleanup()
	return nil
}

//
func (m *LynxiServer) Serve() error {
	os.Remove(m.socket)
	sock, err := net.Listen("unix", m.socket)
	if err != nil {
		return err
	}

	pluginapi.RegisterDevicePluginServer(m.server, *m.plugin)

	go func() {
		lastCrashTime := time.Now()
		restartCount := 0
		for {
			log.Printf("Starting GRPC server for '%s'", m.resourceName)
			err := m.server.Serve(sock)
			if err == nil {
				break
			}

			log.Printf("GRPC server for '%s' crashed with error: %v", m.resourceName, err)

			// restart if it has not been too often
			// i.e. if server has crashed more than 5 times and it didn't last more than one hour each time
			if restartCount > 5 {
				// quit
				log.Fatalf("GRPC server for '%s' has repeatedly crashed recently. Quitting", m.resourceName)
			}
			timeSinceLastCrash := time.Since(lastCrashTime).Seconds()
			lastCrashTime = time.Now()
			if timeSinceLastCrash > 3600 {
				// it has been one hour since the last crash.. reset the count
				// to reflect on the frequency
				restartCount = 1
			} else {
				restartCount++
			}
		}
	}()

	// Wait for server to start by launching a blocking connexion
	conn, err := m.dial(m.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

//
func (m *LynxiServer) Register() error {
	conn, err := m.dial(pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(m.socket),
		ResourceName: m.resourceName,
		Options: &pluginapi.DevicePluginOptions{
			GetPreferredAllocationAvailable: false,
		},
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

func (m *LynxiServer) initialize() {
	m.server = grpc.NewServer([]grpc.ServerOption{}...)
}

func (m *LynxiServer) cleanup() {
	m.server = nil
}

// dial establishes the gRPC communication with the registered device plugin.
func (m *LynxiServer) dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
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
