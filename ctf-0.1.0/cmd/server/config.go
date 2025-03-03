package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/k8s-2025-pschoeppner/ctf/pkg/flagset"
)

type Config struct {
	Flags map[string]string `json:"flags"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	f, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(f, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func newFlagSetFromConfig(cfg Config, flagSet flagset.FlagSet) error {
	for name, value := range cfg.Flags {
		if _, found := flagSet[name]; !found {
			return fmt.Errorf("flag %q not found", name)
		}
		flagSet[name].SetValue(value)
	}
	for name, flag := range flagSet {
		if flag.Value == "" {
			return fmt.Errorf("flag %q not found in config", name)
		}
	}
	return nil
}
