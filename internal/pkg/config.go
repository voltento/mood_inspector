package pkg

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"github.com/voltento/mood_inspector/internal/pkg/notification"
	"gitlab.mobbtech.com/iqbus/iqbus_go_client/errors"
	"log"
)

type (
	TelegramBot struct {
		Id     string `json:"id" mapstructure:"id"`
		Secret string `json:"secret" json:"secret"`
	}

	Config struct {
		TelegramBot   TelegramBot                    `json:"telegram_bot" mapstructure:"telegram_bot"`
		Notifications []notification.NotificationCfg `json:"notifications" mapstructure:"notifications"`
	}
)

func BuildConfig() *Config {
	var path string
	flag.StringVar(&path, "config", "", "path to a json config file")
	flag.Parse()

	if len(path) == 0 {
		log.Fatal("config file path was not provided")
	}

	cfg, err := buildConfigFromFile(path)
	if err != nil {
		log.Fatalf("build config field: %v", err.Error())
	}
	return cfg
}

func buildConfigFromFile(path string) (*Config, error) {
	viper.SetConfigType("json")
	viper.SetConfigFile(path)
	if er := viper.ReadInConfig(); er != nil {
		return nil, errors.Wrap(er, fmt.Sprintf("can not load config from the file %v", path))
	}

	cfg := &Config{}
	if er := viper.Unmarshal(cfg, notification.DecoderConfigOptions()); er != nil {
		return nil, errors.Wrap(er, fmt.Sprintf("can not unmarshal config from the file %v", path))
	}
	return cfg, nil
}
