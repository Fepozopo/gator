package main

import (
	"fmt"
)

// command represents a CLI command.
type command struct {
	name string
	args []string
}

// commands holds all registered command handlers.
type commands struct {
	handlers map[string]func(*state, command) error
}

// register adds a new command handler to the commands map.
func (c *commands) register(name string, handler func(*state, command) error) {
	if c.handlers == nil {
		c.handlers = make(map[string]func(*state, command) error)
	}
	c.handlers[name] = handler
}

// run executes a command if it exists in the handlers map.
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}
