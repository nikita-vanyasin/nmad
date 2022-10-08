package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func newAPIListener(ctx context.Context) *http.Server {
	router := httprouter.New()
	router.GET("/list/:chat_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		chatID := ps.ByName("chat_id")

		nls, err := storage.List(ctx, chatID)
		if err != nil {
			log.Printf("List %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(nls)
		if err != nil {
			log.Printf("Marshal %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, string(response))
	})

	return &http.Server{Addr: "0.0.0.0:" + CONFIG.APIPort, Handler: router}
}
