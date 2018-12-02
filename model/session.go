package model

// Session contains info about session.
type Session struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	Event  string `json:"event"`
}
