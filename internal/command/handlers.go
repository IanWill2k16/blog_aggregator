package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/IanWill2k16/blog_aggregator/internal/config"
	"github.com/IanWill2k16/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(*config.State, Command) error {
	return func(s *config.State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			fmt.Println("user does not exist")
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
}

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
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing variable: agg <time_between_reqs>")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func AddFeed(s *config.State, cmd Command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) < 2 {
		return fmt.Errorf("missing arguments: addfeed <name> <url>")
	}

	url := cmd.Args[1]

	args := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   cmd.Args[0],
		Url:    url,
		UserID: user.ID,
	}
	feedEntry, err := s.Db.CreateFeed(ctx, args)
	if err != nil {
		return err
	}
	followCmd := &Command{
		Args: []string{url},
	}

	if err := Follow(s, *followCmd, user); err != nil {
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
		userName, err := s.Db.GetNameByID(ctx, feeds[i].UserID)
		if err != nil {
			return err
		}
		fmt.Println("")
		fmt.Println(feeds[i].Name)
		fmt.Println(feeds[i].Url)
		fmt.Println(userName)
		fmt.Println("")
		fmt.Println("------------------")
	}
	return nil
}

func Follow(s *config.State, cmd Command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing argument: follow <url>")
	}
	feed, err := s.Db.GetFeedFromURL(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	args := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	res, err := s.Db.CreateFeedFollows(ctx, args)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func Following(s *config.State, cmd Command, user database.User) error {
	ctx := context.Background()
	res, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for i := range res {
		fmt.Println(res[i].FeedName)
	}
	return nil
}

func Unfollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing argument: unfollow <url>")
	}

	args := database.RemoveFeedFollowParams{
		Url:    cmd.Args[0],
		UserID: user.ID,
	}

	err := s.Db.RemoveFeedFollow(context.Background(), args)
	if err != nil {
		return err
	}
	fmt.Println("unfollowed")
	return nil
}
