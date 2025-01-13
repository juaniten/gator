package main

import "errors"

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	commandFunction, ok := c.handlers[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return commandFunction(s, cmd)
}
