package command

import (
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username required: login <username>")
	}
	s.Cfg.SetUser(cmd.Args[0])
	fmt.Println("user has been set")
	return nil
}
