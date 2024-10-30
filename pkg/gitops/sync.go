package gitops

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/tusharbecoding/argocd-clone/pkg/k8s"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Git struct {
		Repo string `yaml:"repo"`
		Branch string `yaml:"branch"`
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
	fmt.Println("Cloning git repo", repo)
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

    for {
        if err := GitPull(); err != nil {
            fmt.Println("Failed to pull:", err)
            return err
        }
        files, err := os.ReadDir("repo/manifests")
        if err != nil {
            return err
        }

        for _, file := range files {
            manifestPath := fmt.Sprintf("repo/manifests/%s", file.Name())
            manifest, err := os.ReadFile(manifestPath)
            if err != nil {
                return err
            }
            fmt.Println("Applying:", manifestPath)
            if err := client.ApplyManifest(string(manifest)); err != nil {
                fmt.Println("Failed to apply manifest:", err)
                return err
            }
        }

        fmt.Println("Sync completed, waiting for next poll interval...")
        time.Sleep(config.Git.PollInterval)
    }
}
