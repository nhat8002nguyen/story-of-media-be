package models

type Story struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	// other fields such as title, author, etc.
}
