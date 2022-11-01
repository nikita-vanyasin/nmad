package main

import (
	"context"
	"log"
	"os"
	"os/signal"
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
