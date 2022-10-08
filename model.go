package main

import "time"

type NomadLocation struct {
	Location
	ChatID   string    `json:"chat_id"`
	Username string    `json:"username"`
	At       time.Time `json:"at"`
}

type Location struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}
