package model

// UserInfo contains info about user.
type UserInfo struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Picture  string `json:"picture"`
}
