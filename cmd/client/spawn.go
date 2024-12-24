package main


func command_spawn(config *Config, args ...string) error {
	err := config.gameState.CommandSpawn(args)
	if err != nil {
		return err
	}
	return nil
}
