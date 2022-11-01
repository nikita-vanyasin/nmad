package main

import "time"

type NomadLocation struct {
	Location `bson:"inline"`
	ChatID   string    `json:"chat_id" bson:"chat_id"`
	Username string    `json:"username" bson:"username"`
	At       time.Time `json:"at" bson:"at"`
}

type Location struct {
	Country string  `json:"country" bson:"country"`
	City    string  `json:"city" bson:"city"`
	Lat     float64 `json:"lat" bson:"lat"`
	Lng     float64 `json:"lng" bson:"lng"`
}
