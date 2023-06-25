package main

import (
	"errors"
	"io/fs"

	"github.com/BurntSushi/toml"
	gap "github.com/muesli/go-app-paths"
)

// config represents the toml configuration file.
type config struct {
	Display string
	Start   string
	End     string
}

// loadConfig loads the config file.
func loadConfig() (*config, error) {
	scope := gap.NewScope(gap.User, "initiative")
	configPath, err := scope.ConfigPath("server.toml")
	if err != nil {
		return &config{}, err
	}

	conf := new(config)
	_, err = toml.DecodeFile(configPath, conf)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return &config{}, err
	}
	return conf, nil
}
