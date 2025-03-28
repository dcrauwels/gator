package rssfeed

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

// describes the entire rss feed. with json for unmarshalling purposes
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// describes a single item in an rss feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (r *RSSFeed) fixTitleDesc() {
	// always clean title & desc
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	// sanity check
	if len(r.Channel.Item) == 0 {
		return
	}

	for i := range r.Channel.Item {
		r.Channel.Item[i].Title = html.UnescapeString(r.Channel.Item[i].Title)
		r.Channel.Item[i].Description = html.UnescapeString(r.Channel.Item[i].Description)
	}

}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var r RSSFeed
	// massively copied from godev documentation examples

	// init client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// make request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &r, err
	}
	req.Header.Set("User-Agent", "gator")

	// Do request
	resp, err := client.Do(req)
	if err != nil {
		return &r, err
	}

	// unpack response into []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &r, err
	}

	// unmarshal body
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return &r, err
	}

	// cleanup
	r.fixTitleDesc()

	return &r, nil
}
