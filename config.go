package httpNostr

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	PrivateKey string `env:"NOSTR_KEY"`
}

var Configuration Config

func init() {
	err := env.Parse(&Configuration)
	if err != nil {
		panic(err)
	}
}
