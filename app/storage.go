package main

import (
	"context"
)

type Storage interface {
	Save(ctx context.Context, nl NomadLocation) error
	List(ctx context.Context, chatID string) ([]NomadLocation, error)
	ListByCountry(ctx context.Context, chatID, country string) ([]NomadLocation, error)
	ListByCity(ctx context.Context, chatID, city string) ([]NomadLocation, error)
	Close(ctx context.Context) error
}
