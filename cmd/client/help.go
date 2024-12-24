package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)
func command_help(config *Config, _ ...string) error {
	gamelogic.PrintClientHelp()
	return nil
}
