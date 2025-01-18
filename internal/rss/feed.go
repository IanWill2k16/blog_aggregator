package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	xmlData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	feedData := &RSSFeed{}
	if err := xml.Unmarshal(xmlData, feedData); err != nil {
		return nil, err
	}

	for i := range feedData.Channel.Item {
		feedData.Channel.Item[i].Title = html.UnescapeString(feedData.Channel.Item[i].Title)
		feedData.Channel.Item[i].Link = html.UnescapeString(feedData.Channel.Item[i].Link)
		feedData.Channel.Item[i].Description = html.UnescapeString(feedData.Channel.Item[i].Description)
		feedData.Channel.Item[i].PubDate = html.UnescapeString(feedData.Channel.Item[i].PubDate)
	}

	return feedData, nil
}
