package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/IanWill2k16/blog_aggregator/internal/config"
	"github.com/IanWill2k16/blog_aggregator/internal/database"
	"github.com/IanWill2k16/blog_aggregator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username required: login <username>")
	}
	if _, err := s.Db.GetUser(context.Background(), cmd.Args[0]); err != nil {
		fmt.Println("user does not exist")
		os.Exit(1)
	}
	s.Cfg.SetUser(cmd.Args[0])
	fmt.Println("user has been set")
	return nil
}

func Register(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username required: login <username>")
	}

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	user, err := s.Db.CreateUser(context.Background(), args)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			fmt.Printf("user already exists\n")
			os.Exit(1)
		}
		return fmt.Errorf("error creating user: %v", err)
	}
	fmt.Printf("Created user: %+v\n", user)
	s.Cfg.SetUser(cmd.Args[0])
	fmt.Println("user has been set")
	return nil
}

func Reset(s *config.State, cmd Command) error {
	if err := s.Db.Reset(context.Background()); err != nil {
		fmt.Println("could not reset database")
		return err
	}
	fmt.Println("database has been reset")
	return nil
}

func GetUsers(s *config.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for i := range users {
		name := users[i].Name
		if name == s.Cfg.CurrentUserName {
			name += " (current)"
		}
		fmt.Printf("* %v\n", name)
	}
	return nil
}

func Agg(s *config.State, cmd Command) error {
	url := "https://www.wagslane.dev/index.xml"
	res, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func AddFeed(s *config.State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("missing arguments: addfeed <name> <url>")
	}
	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user_id not found")
	}

	args := database.CreateFeedParams{
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	}
	feedEntry, err := s.Db.CreateFeed(context.Background(), args)
	if err != nil {
		return err
	}
	fmt.Println(feedEntry)

	return nil
}

func Feeds(s *config.State, cmd Command) error {
	ctx := context.Background()
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	fmt.Println("------------------")
	for i := range feeds {
		user_name, err := s.Db.GetNameByID(ctx, feeds[i].UserID)
		if err != nil {
			return err
		}
		fmt.Println("")
		fmt.Println(feeds[i].Name)
		fmt.Println(feeds[i].Url)
		fmt.Println(user_name)
		fmt.Println("")
		fmt.Println("------------------")
	}
	return nil
}
