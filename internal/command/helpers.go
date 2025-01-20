package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IanWill2k16/blog_aggregator/internal/config"
	"github.com/IanWill2k16/blog_aggregator/internal/database"
	"github.com/IanWill2k16/blog_aggregator/internal/rss"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
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

	for _, item := range feedData.Channel.Item {
		postArgs := database.CreatePostParams{
			ID:     uuid.New(),
			Title:  item.Title,
			Url:    item.Link,
			FeedID: nextFeed.ID,
		}

		if item.Description == "" {
			postArgs.Description = sql.NullString{Valid: false}
		} else {
			postArgs.Description = sql.NullString{String: item.Description, Valid: true}
		}

		parsedTime, err := dateparse.ParseAny(item.PubDate)
		if err != nil {
			postArgs.PublishedAt = sql.NullTime{Valid: false}
		} else {
			postArgs.PublishedAt = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}

		_, err = s.Db.CreatePost(ctx, postArgs)
		var pgErr *pgconn.PgError
		if err != nil {
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				continue
			}
			fmt.Printf("error saving post with URL %s: %v\n", postArgs.Url, err)
			continue
		}
	}
	s.Db.MarkFeedFetched(ctx, nextFeed.ID)

	return nil
}
