package shared

import "time"

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// RssHeadline represents a news headline from an RSS feed.
type RssHeadline struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	PublishedAt string `json:"publishedAt"`
	Source      string `json:"source"`
}
