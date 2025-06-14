package main

import (
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

// FetchRSSFeeds fetches RSS feeds for all users
func FetchRSSFeeds() {
	log.Println("Starting RSS feed fetch...")

	rows, err := db.Query("SELECT id, url FROM feeds")
	if err != nil {
		log.Printf("Error querying feeds: %v", err)
		return
	}
	defer rows.Close()

	fp := gofeed.NewParser()

	for rows.Next() {
		var feedID int
		var feedURL string
		err := rows.Scan(&feedID, &feedURL)
		if err != nil {
			log.Printf("Error scanning feed row: %v", err)
			continue
		}

		log.Printf("Fetching feed: %s", feedURL)
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			log.Printf("Error parsing feed %s: %v", feedURL, err)
			continue
		}

		// Update last_fetched timestamp
		_, err = db.Exec("UPDATE feeds SET last_fetched = NOW() WHERE id = $1", feedID)
		if err != nil {
			log.Printf("Error updating last_fetched for feed %d: %v", feedID, err)
		}

		// Insert new posts
		for _, item := range feed.Items {
			var publishedAt *time.Time
			if item.PublishedParsed != nil {
				publishedAt = item.PublishedParsed
			}

			// Check if post already exists
			var existingID int
			err := db.QueryRow("SELECT id FROM posts WHERE url = $1", item.Link).Scan(&existingID)
			if err == nil {
				// Post already exists, skip
				continue
			}

			// Insert new post
			_, err = db.Exec(
				"INSERT INTO posts (feed_id, title, url, published_at) VALUES ($1, $2, $3, $4)",
				feedID, item.Title, item.Link, publishedAt,
			)
			if err != nil {
				log.Printf("Error inserting post: %v", err)
			}
		}

		log.Printf("Processed %d items from feed %s", len(feed.Items), feedURL)
	}

	log.Println("RSS feed fetch completed")
}

// StartRSSFetcher starts a background goroutine to fetch RSS feeds periodically
func StartRSSFetcher() {
	go func() {
		// Initial fetch
		FetchRSSFeeds()

		// Fetch every 30 minutes
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			FetchRSSFeeds()
		}
	}()
}

