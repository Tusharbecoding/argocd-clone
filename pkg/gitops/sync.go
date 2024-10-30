package gitops

import (
	"os"
	"time"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
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

