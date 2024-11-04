package gitops

import (
	"log"
	"time"

	"github.com/tusharbecoding/argocd-clone/pkg/k8s"
)

// Poller continuously polls the Git repository and applies any changes
type Poller struct {
    config *Config
    client *k8s.K8sClient
}

// NewPoller initializes a new Poller instance
func NewPoller(config *Config, client *k8s.K8sClient) *Poller {
    return &Poller{config: config, client: client}
}

// Start initiates the polling loop
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
