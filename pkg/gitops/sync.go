package gitops

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tusharbecoding/argocd-clone/pkg/k8s"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Git struct {
		Repo         string        `yaml:"repo"`
		Branch       string        `yaml:"branch"`
		PollInterval time.Duration `yaml:"pollInterval"`
	} `yaml:"git"`
	Kubernetes struct {
		Kubeconfig string `yaml:"kubeconfig"`
	} `yaml:"kubernetes"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func GitClone(repo string) error {
	fmt.Println("Cloning git repo:", repo)
	if err := exec.Command("git", "clone", repo, "repo").Run(); err != nil {
		return err
	}
	return nil
}

func GitPull() error {
	cmd := exec.Command("git", "-C", "repo", "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Sync(config *Config, client *k8s.K8sClient) error {
	if _, err := os.Stat("repo"); os.IsNotExist(err) {
		if err := GitClone(config.Git.Repo); err != nil {
			return err
		}
	}

	if err := GitPull(); err != nil {
		fmt.Println("Failed to pull:", err)
		return err
	}

	err := filepath.Walk("repo/manifests", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			fmt.Println("Applying:", path)
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if applyErr := client.ApplyManifest(string(content)); applyErr != nil {
				log.Printf("Failed to apply manifest %s: %v\n", path, applyErr)
				return applyErr
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("Sync completed, waiting for next poll interval...")
	return nil
}


