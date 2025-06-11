package main

import "errors"

type command struct {
	name      string
	arguments []string
}

type commands struct {
	available_commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.available_commands[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.available_commands[name] = f
}
