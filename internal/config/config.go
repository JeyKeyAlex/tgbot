package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
)

func NewConfig() (*Configuration, error) {
	var envFiles []string
	if _, err := os.Stat(".env"); err == nil {
		log.Println("found .env file, adding it to env config files list")
		envFiles = append(envFiles, ".env")
	}

	if len(envFiles) > 0 {
		err := godotenv.Overload(envFiles...)
		if err != nil {
			return nil, errors.Wrapf(err, "error while opening env config: %s", err)
		}
	}

	cfg := &Configuration{}
	ctx := context.Background()

	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "error while loading config")
	}
	return cfg, nil
}

type (
	Configuration struct {
		Telegram Telegram `env:",prefix=TELEGRAM_"`
	}

	Telegram struct {
		ApiToken string `env:"APITOKEN,required"`
	}
)
