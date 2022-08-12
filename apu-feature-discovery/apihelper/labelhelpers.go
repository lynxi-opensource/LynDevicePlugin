package apihelper

import (
	"path"
	"strings"
	"log"
	"context"
	"k8s.io/apimachinery/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"encoding/json"
)

const (
	FeatureLabelNs = "feature.node.kubernetes.io"
	featureLabelAnnotation = "feature-labels"
	AnnotationNsBase = "nfd.node.kubernetes.io"
)

type Labels map[string]string

type LabelHelpers struct {
	nodeName 	 string
	annotationNs string
	apihelper    APIHelpers
}

func NewHelper(name string) *LabelHelpers {
	return &LabelHelpers{
		nodeName:       name,
	}
}

func (m *LabelHelpers) Init() error {
	m.annotationNs = AnnotationNsBase

	kubeconfig, err := GetKubeconfig("")
	if err != nil {
		return err
	}
	m.apihelper = K8sHelpers{Kubeconfig: kubeconfig}
	return nil
}

//更新node标签
func (m *LabelHelpers) UpdateNodeLabels(labels Labels) error {

	cli, err := m.apihelper.GetClient()
	if err != nil {
		log.Printf("GetClient error ")
		return err
	}

	ctx := context.TODO()

	// Get the worker node object
	node, err := m.apihelper.GetNode(cli, m.nodeName)
	if err != nil {
		log.Printf("GetNode error %v", err)
		return err
	}

	labelsNew := node.Labels
	for key, value := range labels {
		labelsNew[key] = value
	}
	
	patchData := map[string]interface{}{"metadata": map[string]map[string]string{"labels": labelsNew}}
	playLoadBytes, _ := json.Marshal(patchData)
	_, err = cli.CoreV1().Nodes().Patch(ctx, m.nodeName, types.StrategicMergePatchType, playLoadBytes, metav1.PatchOptions{})
	if err != nil {
		log.Printf("[ModifyNodeLabels] %v node Patch fail %v\n", m.nodeName, err)
		log.Printf("error：%s", err)
	}

	// log.Printf("modify node %s label %s to %s", nodeName, labelKey, labelNewValue)
	return err
}


func (m *LabelHelpers) annotationName(name string) string {
	return path.Join(m.annotationNs, name)
}

// stringToNsNames is a helper for converting a string of comma-separated names
// into a slice of fully namespaced names
func stringToNsNames(cslist, ns string) []string {
	var names []string
	if cslist != "" {
		names = strings.Split(cslist, ",")
		for i, name := range names {
			// Expect that names may omit the ns part
			names[i] = addNs(name, ns)
		}
	}
	return names
}

// addNs adds a namespace if one isn't already found from src string
func addNs(src string, nsToAdd string) string {
	if strings.Contains(src, "/") {
		return src
	}
	return path.Join(nsToAdd, src)
}

// createPatches is a generic helper that returns json patch operations to perform
func createPatches(removeKeys []string, oldItems map[string]string, newItems map[string]string, jsonPath string) []JsonPatch {
	patches := []JsonPatch{}

	// Determine items to remove
	for _, key := range removeKeys {
		if _, ok := oldItems[key]; ok {
			if _, ok := newItems[key]; !ok {
				patches = append(patches, NewJsonPatch("remove", jsonPath, key, ""))
			}
		}
	}

	// Determine items to add or replace
	for key, newVal := range newItems {
		if oldVal, ok := oldItems[key]; ok {
			if newVal != oldVal {
				patches = append(patches, NewJsonPatch("replace", jsonPath, key, newVal))
			}
		} else {
			patches = append(patches, NewJsonPatch("add", jsonPath, key, newVal))
		}
	}

	return patches
}