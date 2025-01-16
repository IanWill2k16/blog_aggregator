package command

import (
	"fmt"
)

func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required: login <username>")
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Println("user has been set")
	return nil
}
