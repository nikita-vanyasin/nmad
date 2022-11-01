package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	APIPort            string `env:"NMAD_API_PORT"`
	TelegramAPIToken   string `env:"TELEGRAM_API_TOKEN"`
	GeoNamesAPILogin   string `env:"GEONAMES_API_LOGIN_NAME"`
	MapServiceEndpoint string `env:"MAP_SERVICE_ENDPOINT"`

	ArangoDBEndpoint string `env:"ARANGODB_ENDPOINT" envDefault:""`
	ArangoDBUser     string `env:"ARANGODB_USER" envDefault:""`
	ArangoDBPassword string `env:"ARANGODB_PASSWORD" envDefault:""`
	ArangoDBCA       string `env:"ARANGODB_CA" envDefault:""`

	MongoDBConnectURL string `env:"MONGODB_CONNECT_URL" envDefault:""`
}

var CONFIG Config

func init() {
	err := env.Parse(&CONFIG, env.Options{RequiredIfNoDef: true})
	if err != nil {
		panic(err)
	}
}
