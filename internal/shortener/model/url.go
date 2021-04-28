package model

import "time"

type URL struct {
	ID int `json:"id,omitempty"`

	UserID       int64     `json:"owner_id"`
	ShortenedURL string    `json:"shortened"`
	RealURL      string    `json:"real"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
