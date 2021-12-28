package model

import (
	"fmt"
	"time"
)

type URL struct {
	ID int `json:"id,omitempty"`

	UserID       int64     `json:"owner_id"`
	ShortenedURL string    `json:"shortened"`
	RealURL      string    `json:"real"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

func (u *URL) String() string {
	return fmt.Sprintf("{ Owner: %d, Shortened: %s, RealURL: %s, CreatedAt: %v, UpdatedAt: %v }", u.UserID, u.ShortenedURL, u.RealURL, u.CreatedAt.UTC().String(), u.UpdatedAt)
}
