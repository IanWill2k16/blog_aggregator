package command

import (
	"fmt"

	"github.com/IanWill2k16/blog_aggregator/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

type State struct {
	Cfg *config.Config
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandMap[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.CommandMap[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	return handler(s, cmd)
}
