package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func newHTTPListener(ctx context.Context, distDir string) *http.Server {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.Use(static.Serve("/", static.LocalFile(distDir, true)))

	api := r.Group("/api/v1")
	api.GET("/list/:chat_id", func(c *gin.Context) {
		chatID := c.Param("chat_id")

		nls, err := storage.List(ctx, chatID)
		if err != nil {
			log.Printf("List: %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		type respItem struct {
			NomadLocation
			ProfileURL string `json:"profile_url"`
		}
		respItems := make([]respItem, 0)
		for _, l := range nls {
			respItems = append(respItems, respItem{
				NomadLocation: l,
				ProfileURL:    fmt.Sprintf("https://t.me/%s", l.Username),
			})
		}

		response, err := json.Marshal(respItems)
		if err != nil {
			log.Printf("Marshal %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.String(http.StatusOK, string(response))
	})

	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(distDir, "index.html"))
	})

	return &http.Server{Addr: "0.0.0.0:" + CONFIG.APIPort, Handler: r}
}
