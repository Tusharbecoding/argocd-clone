package k8s

import (
	"bytes"
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	cached "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	Clientset     *kubernetes.Clientset
	DynamicClient dynamic.Interface
	Config        *rest.Config
	RESTMapper    meta.RESTMapper
}

func NewK8sClient(kubeconfig string) (*K8sClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	discoveryClient := cached.NewMemCacheClient(clientset.Discovery())
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)

	return &K8sClient{
		Clientset:     clientset,
		DynamicClient: dynClient,
		Config:        config,
		RESTMapper:    mapper,
	}, nil
}

func (k *K8sClient) ApplyManifest(manifest string) error {
	fmt.Println("Applying manifest:", manifest)

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(manifest)), 4096)

	obj := &unstructured.Unstructured{}
	if err := decoder.Decode(obj); err != nil {
		return fmt.Errorf("failed to decode manifest: %v", err)
	}

	gvk := obj.GroupVersionKind()
	mapping, err := k.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("failed to get REST mapping: %v", err)
	}

	namespace := obj.GetNamespace()
	if namespace == "" {
		namespace = "default"
		obj.SetNamespace(namespace)
	}

	resourceClient := k.DynamicClient.Resource(mapping.Resource).Namespace(namespace)
	_, err = resourceClient.Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Printf("Resource %s already exists, updating...\n", obj.GetName())
			_, err = resourceClient.Update(context.TODO(), obj, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("failed to update existing resource: %v", err)
			}
		} else {
			return fmt.Errorf("failed to create resource: %v", err)
		}
	}

	fmt.Println("Manifest applied successfully")
	return nil
}
