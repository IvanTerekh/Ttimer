package model

type Session struct {
	Name   string `json:"name"`
	UserId string
	Event  string `json:"event"`
}
