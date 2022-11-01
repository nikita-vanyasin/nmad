package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	tele "gopkg.in/telebot.v3"
)

const commandPrefix = "nmad "

var (
	storage Storage
	geoInfo GeoInfo
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error
	storage, err = NewMongoDBStorage()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer storage.Close(ctx)

	geoInfo, err = NewGeoInfo()
	if err != nil {
		log.Fatal(err)
		return
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	b := newTelegramBot(ctx)
	go b.Start()

	apiServer := newAPIListener(ctx)
	go func() {
		err := apiServer.ListenAndServe()
		if err != nil && err != context.Canceled {
			log.Printf("API server stopped with error %s", err)
		}
	}()
	log.Printf("Started!")

	<-stop
	log.Printf("got SIGINT. Terminating...")
	b.Stop()
	cancel()
}

func handleCommand(ctx context.Context, cmd string, chat *tele.Chat, sender *tele.User) (string, error) {
	m := map[string]HandlerFunc{
		"city set ":    handleSetCity,
		"city get ":    handleGetCity,
		"country get ": handleGetCountry,
		"list":         handleList,
		"map":          handleMap,
	}

	usageInfo := "Available commands:"
	for p := range m {
		usageInfo += fmt.Sprintf("\n%s%s", commandPrefix, strings.TrimSpace(p))
	}
	m["help"] = func(ctx context.Context, i []string, chat *tele.Chat, user *tele.User) (string, error) {
		return usageInfo, nil
	}

	for p, fn := range m {
		if rest := strings.TrimPrefix(cmd, p); rest != cmd {
			args := strings.Split(rest, " ")
			args = filterOutEmptyStrings(args)
			return fn(ctx, args, chat, sender)
		}
	}

	return fmt.Sprintf("Unknown command. %s", usageInfo), nil
}

func filterOutEmptyStrings(strs []string) []string {
	var result []string
	for _, s := range strs {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
