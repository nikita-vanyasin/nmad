package main

import (
	"context"
	"log"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func newTelegramBot(ctx context.Context) *tele.Bot {
	pref := tele.Settings{
		Token:  CONFIG.TelegramAPIToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil
	}
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

	return b
}
