package gitops

import (
	"log"
	"time"

	"github.com/tusharbecoding/argocd-clone/pkg/k8s"
)

type Poller struct {
    config *Config
    client *k8s.K8sClient
}

func NewPoller(config *Config, client *k8s.K8sClient) *Poller {
    return &Poller{config: config, client: client}
}


func (p *Poller) Start() {
    ticker := time.NewTicker(p.config.Git.PollInterval)
    defer ticker.Stop()

    for range ticker.C {
        log.Println("Polling for changes...")
        err := Sync(p.config, p.client)
        if err != nil {
            log.Printf("Error during sync: %v", err)
        }
    }
}
