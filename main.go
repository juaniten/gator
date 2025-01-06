package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/juaniten/gator/internal/config"
)

func main() {
	// Read JSON config file and load into state
	configuration, err := config.Read()
	if err != nil {
		log.Fatalf("error reading gator configuration: %v", err)
	}
	state_pointer := &state{
		config: &configuration,
	}

	// Initialize commands
	comm := commands{
		handlers: map[string]func(*state, command) error{},
	}
	comm.register("login", handlerLogin)

	// Process arguments
	args := os.Args
	if len(args) < 3 {
		log.Fatal("command needs an argument")
	}
	// Execute command
	newCommand := command{
		name:      args[1],
		arguments: args[2:],
	}
	err = comm.run(state_pointer, newCommand)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
	}
}

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if commandFunction, ok := c.handlers[cmd.name]; ok {
		commandFunction(s, cmd)
		return nil
	}
	fmt.Println("command does not exist")
	return nil
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		return errors.New("command called with no arguments")
	}
	username := cmd.arguments[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("Set username '%s'\n", username)
	return nil
}
