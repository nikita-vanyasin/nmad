package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func newAPIListener(ctx context.Context) *http.Server {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/api/v1/list/:chat_id", func(c *gin.Context) {
		chatID := c.Param("chat_id")

		nls, err := storage.List(ctx, chatID)
		if err != nil {
			log.Printf("List: %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(nls)
		if err != nil {
			log.Printf("Marshal %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.String(http.StatusOK, string(response))
	})

	return &http.Server{Addr: "0.0.0.0:" + CONFIG.APIPort, Handler: r}
}
