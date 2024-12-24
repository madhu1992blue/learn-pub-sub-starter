package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)
func command_quit(config *Config, _ ...string) error {
	gamelogic.PrintQuit()
	return nil
}
