package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
)

// GetUserFromToken extracts username from JWT token
func GetUserFromToken(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(c.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Username, nil
}

// GetUserID gets user ID from username
func GetUserID(username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	return userID, err
}

// GetFeeds returns all feeds for a user
func GetFeeds(w http.ResponseWriter, r *http.Request) {
	username, err := GetUserFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := GetUserID(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, name, url, last_fetched FROM feeds WHERE user_id = $1", userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var feed Feed
		var lastFetched sql.NullString
		err := rows.Scan(&feed.ID, &feed.Name, &feed.URL, &lastFetched)
		if err != nil {
			continue
		}
		feed.UserID = userID
		if lastFetched.Valid {
			feed.LastFetched = lastFetched.String
		}
		feeds = append(feeds, feed)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

// AddFeed adds a new feed for a user
func AddFeed(w http.ResponseWriter, r *http.Request) {
	username, err := GetUserFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := GetUserID(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var feed Feed
	err = json.NewDecoder(r.Body).Decode(&feed)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO feeds (user_id, name, url) VALUES ($1, $2, $3)", userID, feed.Name, feed.URL)
	if err != nil {
		log.Printf("Error inserting feed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Trigger immediate RSS fetch for the new feed
	go FetchRSSFeeds()

	w.WriteHeader(http.StatusCreated)
}

// DeleteFeed deletes a feed for a user
func DeleteFeed(w http.ResponseWriter, r *http.Request) {
	username, err := GetUserFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := GetUserID(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	feedIDStr := chi.URLParam(r, "feedID")
	feedID, err := strconv.Atoi(feedIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM feeds WHERE id = $1 AND user_id = $2", feedID, userID)
	if err != nil {
		log.Printf("Error deleting feed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetPosts returns all posts for a user's feeds
func GetPosts(w http.ResponseWriter, r *http.Request) {
	username, err := GetUserFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := GetUserID(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := `
		SELECT p.id, p.feed_id, p.title, p.url, p.published_at 
		FROM posts p 
		JOIN feeds f ON p.feed_id = f.id 
		WHERE f.user_id = $1 
		ORDER BY p.published_at DESC 
		LIMIT 50
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var publishedAt sql.NullString
		err := rows.Scan(&post.ID, &post.FeedID, &post.Title, &post.URL, &publishedAt)
		if err != nil {
			continue
		}
		if publishedAt.Valid {
			post.PublishedAt = publishedAt.String
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

