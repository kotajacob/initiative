package main

import (
	"errors"
	"io/fs"

	"github.com/BurntSushi/toml"
	gap "github.com/muesli/go-app-paths"
)

// config represents the toml configuration file.
type config struct {
	// Address is for the initiative server.
	Address string

	// Party is an ordered list of default party members.
	Party []string
}

// loadConfig loads the config file.
func loadConfig() (*config, error) {
	scope := gap.NewScope(gap.User, "initiative")
	configPath, err := scope.ConfigPath("client.toml")
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
