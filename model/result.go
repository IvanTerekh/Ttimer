// Package model contains structures definitions.
package model

// Result contains info about solve result. Centis stands for time in centiseconds.
type Result struct {
	Centis   int    `json:"centis"`
	Scramble string `json:"scramble"`
	Penalty  bool   `json:"penalty"`
	Datetime string `json:"datetime"`
}
