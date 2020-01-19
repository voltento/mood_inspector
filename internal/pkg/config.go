package pkg

import (
	"flag"
	"github.com/spf13/viper"
	"log"
)

type (
	TelegramBot struct {
		Id     string `json:"id" mapstructure:"id"`
		Secret string `json:"secret" json:"secret"`
	}

	Config struct {
		TelegramBot TelegramBot `json:"telegram_bot" mapstructure:"telegram_bot"`
	}
)

func BuildConfig() *Config {
	cfg := &Config{}

	var path string
	flag.StringVar(&path, "config", "", "path to a json config file")
	flag.Parse()

	if len(path) == 0 {
		log.Fatal("config file path was not provided")
	} else {
		viper.SetConfigType("json")
		viper.SetConfigFile(path)
		if er := viper.ReadInConfig(); er != nil {
			log.Fatalf("can not load config from the file %v: %v", path, er.Error())
		}
		if er := viper.Unmarshal(cfg); er != nil {
			log.Fatalf("can not load config from the file %v: %v", path, er.Error())
		}
	}
	return cfg
}
