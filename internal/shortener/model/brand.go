package model

type Brand struct {
	ID int `json:"id"`

	UserID int    `json:"owner_id"`
	Brand  string `json:"brand"`
}
