package command

import (
	"fmt"

	"github.com/IanWill2k16/blog_aggregator/internal/config"
)

type command struct {
	name string
	args []string
}

type Commands struct {
	commandMap map[string]func(*State, command) error
}

type State struct {
	cfg *config.Config
}

func (c *Commands) register(name string, f func(*State, command) error) {
	c.commandMap[name] = f
}

func (c *Commands) run(s *State, cmd command) error {
	handler, ok := c.commandMap[cmd.name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	return handler(s, cmd)
}
