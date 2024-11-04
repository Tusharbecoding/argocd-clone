// pkg/gitops/poller.go
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

// NewPoller creates a new Poller instance
func NewPoller(config *Config, client *k8s.K8sClient) *Poller {
    return &Poller{config: config, client: client}
}

// Start initiates the polling mechanism
func (p *Poller) Start() {
    ticker := time.NewTicker(p.config.Git.PollInterval)
    defer ticker.Stop()

    for range ticker.C {
        log.Println("Polling for changes...")
        err := Sync(p.config, p.client)
        if err != nil {
            log.Printf("Error syncing: %v", err)
        }
    }
}
