package main

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Don't expose password in JSON
}

type Feed struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	LastFetched string `json:"last_fetched"`
}

type Post struct {
	ID          int    `json:"id"`
	FeedID      int    `json:"feed_id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	PublishedAt string `json:"published_at"`
}

