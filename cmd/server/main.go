package main

import (
	"log"

	"github.com/tusharbecoding/argocd-clone/pkg/gitops"
	"github.com/tusharbecoding/argocd-clone/pkg/k8s"
)

func main() {
    config, err := gitops.LoadConfig("configs/config.yaml")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    client, err := k8s.NewK8sClient(config.Kubernetes.Kubeconfig)
    if err != nil {
        log.Fatal("Failed to create Kubernetes client:", err)
    }

    poller := gitops.NewPoller(config, client)
    poller.Start()
}
