package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/IanWill2k16/blog_aggregator/internal/command"
	"github.com/IanWill2k16/blog_aggregator/internal/config"
	"github.com/IanWill2k16/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	state := config.State{}

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	state.Cfg = &cfg

	db, err := sql.Open("postgres", state.Cfg.DBUrl)
	if err != nil {
		fmt.Println("error connection to database: %w", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	state.Db = dbQueries

	commands := command.Commands{
		CommandMap: make(map[string]func(*config.State, command.Command) error),
	}

	commands.Register("login", command.HandlerLogin)
	commands.Register("register", command.Register)
	commands.Register("reset", command.Reset)
	commands.Register("users", command.GetUsers)
	commands.Register("agg", command.Agg)
	commands.Register("addfeed", command.AddFeed)
	commands.Register("feeds", command.Feeds)

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
