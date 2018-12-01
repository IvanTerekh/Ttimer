package model

// Session contains info about session.
type Session struct {
	Name   string `json:"name"`
	UserId string `json:"user_id"`
	Event  string `json:"event"`
}
