package main

import (
	"context"
	"log"
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

	m := map[string]HandlerFunc{
		"set":  handleSet,
		"get":  handleGet,
		"list": handleList,
		"map":  handleMap,
	}
	for cmd, handler := range m {
		handler := handler
		b.Handle("/"+cmd, func(c tele.Context) error {
			chat := c.Chat()
			if chat == nil {
				return c.Send("bot can operate only on chats")
			}
			user := c.Sender()
			if user == nil {
				return c.Send("that was unexpected!")
			}

			response, err := handler(ctx, c.Args(), chat, user)
			if err != nil {
				log.Printf("ERR %s", err.Error())
				return c.Send("server error")
			}
			return c.Send(response)
		})
	}

	return b
}
