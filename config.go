package main

import (
	"fmt"
	"os"
)

type Config struct {
	TelegramAPIToken string
	GeoNamesAPILogin string
	ArangoDBEndpoint string
	ArangoDBUser     string
	ArangoDBPassword string
	ArangoDBCA       string
}

var CONFIG Config

func init() {
	CONFIG.TelegramAPIToken = mustGetEnv("TELEGRAM_API_TOKEN")
	CONFIG.GeoNamesAPILogin = mustGetEnv("GEONAMES_API_LOGIN_NAME")
	CONFIG.ArangoDBEndpoint = mustGetEnv("ARANGODB_ENDPOINT")
	CONFIG.ArangoDBUser = mustGetEnv("ARANGODB_USER")
	CONFIG.ArangoDBPassword = mustGetEnv("ARANGODB_PASSWORD")
	CONFIG.ArangoDBCA = mustGetEnv("ARANGODB_CA")
}

func mustGetEnv(n string) string {
	v := os.Getenv(n)
	if v == "" {
		panic(fmt.Sprintf("env variable %s not set", n))
	}
	return v
}
