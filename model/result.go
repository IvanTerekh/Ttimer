package model

type Result struct {
	Centis   int    `json:"centis"`
	Scramble string `json:"scramble"`
	Penalty  bool   `json:"penalty"`
	Datetime string `json:"datetime"`
}
