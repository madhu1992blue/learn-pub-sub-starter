package main
import "fmt"

func command_spam(_ *Config, _ ...string) error {
	fmt.Println("Spamming not allowed yet!")
	return nil
}
