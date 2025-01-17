package main

import (
	"fmt"
	"os"

	"github.com/IanWill2k16/blog_aggregator/internal/command"
	"github.com/IanWill2k16/blog_aggregator/internal/config"
)

func main() {
	state := command.State{}

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	state.Cfg = &cfg

	commands := command.Commands{
		CommandMap: make(map[string]func(*command.State, command.Command) error),
	}

	commands.Register("login", command.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("at least one argument required")
		os.Exit(1)
	}
	cmd := command.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := commands.Run(&state, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
