package main

import (
	"fmt"
)

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
