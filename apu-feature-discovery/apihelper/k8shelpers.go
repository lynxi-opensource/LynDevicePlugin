/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apihelper

import (
	"context"
	"encoding/json"

	api "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Implements APIHelpers
type K8sHelpers struct {
	Kubeconfig *restclient.Config
}

func (h K8sHelpers) GetClient() (*k8sclient.Clientset, error) {
	if h.Kubeconfig == nil {
		cfg, e := GetKubeconfig("")
		if e != nil {
			return nil, e
		}
		h.Kubeconfig = cfg
	}

	clientset, err := k8sclient.NewForConfig(h.Kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func (h K8sHelpers) GetNode(cli *k8sclient.Clientset, nodeName string) (*api.Node, error) {
	// Get the node object using node name
	node, err := cli.CoreV1().Nodes().Get(context.TODO(), nodeName, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (h K8sHelpers) GetNodes(cli *k8sclient.Clientset) (*api.NodeList, error) {
	return cli.CoreV1().Nodes().List(context.TODO(), meta_v1.ListOptions{})
}

func (h K8sHelpers) UpdateNode(c *k8sclient.Clientset, n *api.Node) error {
	// Send the updated node to the apiserver.
	_, err := c.CoreV1().Nodes().Update(context.TODO(), n, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (h K8sHelpers) PatchNode(c *k8sclient.Clientset, nodeName string, patches []JsonPatch) error {
	if len(patches) > 0 {
		data, err := json.Marshal(patches)
		if err == nil {
			_, err = c.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.JSONPatchType, data, meta_v1.PatchOptions{})
		}
		return err
	}
	return nil
}

func (h K8sHelpers) PatchNodeStatus(c *k8sclient.Clientset, nodeName string, patches []JsonPatch) error {
	if len(patches) > 0 {
		data, err := json.Marshal(patches)
		if err == nil {
			_, err = c.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.JSONPatchType, data, meta_v1.PatchOptions{}, "status")
		}
		return err
	}
	return nil

}

func (h K8sHelpers) GetPod(cli *k8sclient.Clientset, namespace string, podName string) (*api.Pod, error) {
	// Get the node object using pod name
	pod, err := cli.CoreV1().Pods(namespace).Get(context.TODO(), podName, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func GetKubeconfig(path string) (*restclient.Config, error) {
	if path == "" {
		return restclient.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", path)
}
