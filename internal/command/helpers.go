package command

import (
	"context"
	"fmt"

	"github.com/IanWill2k16/blog_aggregator/internal/config"
	"github.com/IanWill2k16/blog_aggregator/internal/rss"
)

func scrapeFeeds(s *config.State) error {
	ctx := context.Background()
	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	feedData, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	s.Db.MarkFeedFetched(ctx, nextFeed.ID)

	fmt.Printf("Feed: %v\n", feedData.Channel.Title)
	fmt.Printf("Description: %v\n", feedData.Channel.Description)
	for _, item := range feedData.Channel.Item {
		fmt.Printf("- %v\n", item.Title)
	}
	fmt.Println()

	return nil
}
