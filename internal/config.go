package internal

import (
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

const (
	OutputConsole  = "console"
	OutputMarkdown = "markdown"
)

type Config struct {
	Output     string
	Kubeconfig string
}

func (c Config) Validate() error {
	switch c.Output {
	case OutputConsole, OutputMarkdown:
	default:
		return fmt.Errorf("invalid value for output parameter. Supported values: [%s, %s]", OutputConsole, OutputMarkdown)
	}

	if c.Kubeconfig == "" {
		return fmt.Errorf("missing kubeconfig")
	}

	return nil
}

func Load() (Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.Output, "output", OutputConsole, fmt.Sprintf("output type. Possible values: [%s, %s]", OutputConsole, OutputMarkdown))

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&cfg.Kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&cfg.Kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
