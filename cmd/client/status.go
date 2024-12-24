package main


func command_status(config *Config, _ ...string) error {
	config.gameState.CommandStatus()
	return nil
}
