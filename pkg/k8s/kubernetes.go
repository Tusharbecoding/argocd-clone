package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	Clientset *kubernetes.Clientset
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
    return &K8sClient{Clientset: clientset}, nil
}

func (k *K8sClient) ApplyManifest(manifest string) error {
    fmt.Println("Applying manifest:", manifest)
    return nil
}