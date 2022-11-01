package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
)

type HandlerFunc func(context.Context, []string, *tele.Chat, *tele.User) (string, error)

// lookupCity return city location if city is found. Otherwise bot text response or error
func lookupCity(ctx context.Context, requestedCity string) (*Location, string, error) {
	loc, err := geoInfo.LookupCity(ctx, requestedCity)
	if err != nil {
		return nil, "", errors.WithMessagef(err, "LookupCity %s", requestedCity)
	}
	if loc == nil {
		return nil, "Unknown city :(", nil
	}

	if strings.ToLower(loc.City) != strings.ToLower(requestedCity) {
		return nil, fmt.Sprintf("Unknown city! Did you mean %s?", loc.City), nil

	}
	return loc, "", nil
}

func handleSetCity(ctx context.Context, requestedCity string, chat *tele.Chat, sender *tele.User) (string, error) {
	loc, resp, err := lookupCity(ctx, requestedCity)
	if err != nil {
		return "", err
	}
	if loc == nil {
		return resp, nil
	}

	nl := NomadLocation{
		ChatID:   strconv.FormatInt(chat.ID, 10),
		Username: sender.Username,
		At:       time.Now(),
		Location: *loc,
	}
	err = storage.Save(ctx, nl)
	if err != nil {
		return "", errors.WithMessagef(err, "storage.Save")
	}

	return fmt.Sprintf("City %s set", loc.City), nil
}

func handleSet(ctx context.Context, args []string, chat *tele.Chat, sender *tele.User) (string, error) {
	if len(args) < 2 {
		return "invalid cmd", nil
	}
	switch args[0] {
	case "city":
		return handleSetCity(ctx, strings.Join(args[1:], " "), chat, sender)
	}
	return "invalid cmd", nil
}

func handleGet(ctx context.Context, args []string, chat *tele.Chat, sender *tele.User) (string, error) {
	if len(args) < 2 {
		return "invalid cmd", nil
	}
	restArgs := strings.Join(args[1:], " ")
	switch args[0] {
	case "city":
		return handleGetCity(ctx, restArgs, chat)
	case "country":
		return handleGetCountry(ctx, restArgs, chat)
	}
	return "invalid cmd", nil
}

func handleGetCity(ctx context.Context, requestedCity string, chat *tele.Chat) (string, error) {
	loc, resp, err := lookupCity(ctx, requestedCity)
	if err != nil {
		return "", err
	}
	if loc == nil {
		return resp, nil
	}

	chatID := strconv.FormatInt(chat.ID, 10)
	nls, err := storage.ListByCity(ctx, chatID, loc.City)
	if err != nil {
		return "", errors.WithMessagef(err, "ListByCity %s %s", chatID, loc.City)
	}

	if len(nls) == 0 {
		return fmt.Sprintf("There's no nomads in city %s", loc.City), nil
	}
	var nomadList []string
	for _, nl := range nls {
		nomadList = append(nomadList, fmt.Sprintf("@%s", nl.Username))
	}

	return fmt.Sprintf("Nomads in city %s:\n%s", loc.City, strings.Join(nomadList, "\n")), nil
}

func handleGetCountry(ctx context.Context, requestedCountry string, chat *tele.Chat) (string, error) {
	country, err := geoInfo.LookupCountry(ctx, requestedCountry)
	if err != nil {
		return "", errors.WithMessagef(err, "LookupCountry %s", requestedCountry)
	}
	if country == "" {
		return "Unknown country :(", nil
	}

	if strings.ToLower(country) != strings.ToLower(requestedCountry) {
		return fmt.Sprintf("Unknown country! Did you mean %s?", country), nil
	}

	chatID := strconv.FormatInt(chat.ID, 10)
	nls, err := storage.ListByCountry(ctx, chatID, country)
	if err != nil {
		return "", errors.WithMessagef(err, "ListByCountry %s", country)
	}

	if len(nls) == 0 {
		return fmt.Sprintf("There's no nomads in country %s", country), nil
	}
	var nomadList []string
	for _, nl := range nls {
		nomadList = append(nomadList, fmt.Sprintf("@%s", nl.Username))
	}

	return fmt.Sprintf("Nomads in country %s:\n%s", country, strings.Join(nomadList, "\n")), nil
}

func handleList(ctx context.Context, args []string, chat *tele.Chat, sender *tele.User) (string, error) {
	if len(args) > 0 {
		return "unsupported command", nil
	}

	chatID := strconv.FormatInt(chat.ID, 10)
	nls, err := storage.List(ctx, chatID)
	if err != nil {
		return "", errors.WithMessagef(err, "List")
	}

	if len(nls) == 0 {
		return fmt.Sprintf("Zero nomads in this chat!"), nil
	}
	var nomadList []string
	for _, nl := range nls {
		nomadList = append(nomadList, fmt.Sprintf("@%s", nl.Username))
	}

	return fmt.Sprintf("Nomads in this chat:\n%s", strings.Join(nomadList, "\n")), nil
}

func handleMap(ctx context.Context, args []string, chat *tele.Chat, sender *tele.User) (string, error) {
	if len(args) > 0 {
		return "unsupported command", nil
	}

	chatID := strconv.FormatInt(chat.ID, 10)
	url := fmt.Sprintf("%s/map/%s", CONFIG.MapServiceEndpoint, chatID)
	return fmt.Sprintf("🧭 %s", url), nil
}
