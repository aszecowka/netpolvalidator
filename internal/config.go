package internal

import (
	"flag"
	"fmt"
)

const (
	OutputConsole  = "console"
	OutputMarkdown = "markdown"
)

type Config struct {
	Output string
}

func (c Config) Validate() error {
	switch c.Output {
	case OutputConsole, OutputMarkdown:
	default:
		return fmt.Errorf("invalid value for output parameter. Supported values: [%s, %s]", OutputConsole, OutputMarkdown)
	}

	return nil
}

func Load() (Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.Output, "output", OutputConsole, fmt.Sprintf("output type. Possible values: [%s, %s]", OutputConsole, OutputMarkdown))
	flag.Parse()
	return cfg, nil
}
