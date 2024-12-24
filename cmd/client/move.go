package main


func command_move(config *Config, args ...string) error {
	_, err := config.gameState.CommandMove(args)
	if err != nil {
		return err
	}	
	return nil
}
