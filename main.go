package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const commandPrefix = "nmad "

var (
	storage Storage
	geoInfo GeoInfo
)

func main() {
	var err error
	storage, err = NewArangodbDBStorage()
	if err != nil {
		log.Fatal(err)
		return
	}

	geoInfo, err = NewGeoInfo()
	if err != nil {
		log.Fatal(err)
		return
	}

	pref := tele.Settings{
		Token:  CONFIG.TelegramAPIToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	b.Use(middleware.IgnoreVia())
	b.Use(middleware.Recover(func(err error) {
		log.Printf("ERR %s", err.Error())
	}))

	b.Handle(tele.OnText, func(c tele.Context) error {
		text := c.Text()
		if strings.HasPrefix(text, commandPrefix) {
			chat := c.Chat()
			if chat == nil {
				return c.Send("bot can operate only on chats")
			}

			user := c.Sender()
			if user == nil {
				return c.Send("that was unexpected!")
			}

			response, err := handleCommand(ctx, strings.TrimSpace(strings.TrimPrefix(text, commandPrefix)), chat, user)
			if err != nil {
				log.Printf("ERR %s", err.Error())
				return c.Send("server error")
			}
			return c.Send(response)
		}
		return nil
	})

	go b.Start()
	log.Printf("Started!")

	<-stop
	log.Printf("got SIGINT. Terminating...")
	cancel()
	b.Stop()
}

func handleCommand(ctx context.Context, cmd string, chat *tele.Chat, sender *tele.User) (string, error) {
	m := map[string]HandlerFunc{
		"city set ":    handleSetCity,
		"city get ":    handleGetCity,
		"country get ": handleGetCountry,
		"list":         handleList,
	}

	for p, fn := range m {
		if rest := strings.TrimPrefix(cmd, p); rest != cmd {
			args := strings.Split(rest, " ")
			args = filterOutEmptyStrings(args)
			return fn(ctx, args, chat, sender)
		}
	}

	usageInfo := "Unknown command! Use one of these:"
	for p := range m {
		usageInfo += fmt.Sprintf("\n%s%s", commandPrefix, strings.TrimSpace(p))
	}
	return usageInfo, nil
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
