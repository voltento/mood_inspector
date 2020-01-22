package notification

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"reflect"
	"time"
)

type (
	RandomTimeCfg struct {
		From        time.Time     `json:"from" mapstructure:"from"`
		To          time.Time     `json:"to" mapstructure:"to"`
		Period      time.Duration `json:"period" mapstructure:"period"`
		ExtraPeriod time.Duration `json:"extra_period" mapstructure:"extra_period"`
	}

	NotificationCfg struct {
		Name          string        `json:"name" mapstructure:"name"`
		Message       string        `json:"message" mapstructure:"message"`
		RandomMessage []string      `json:"random_message" mapstructure:"random_message"`
		RandomTime    RandomTimeCfg `json:"random_time" mapstructure:"random_time"`
		CertainTime   []time.Time   `json:"certain_time" mapstructure:"certain_time"`
	}
)

func DecoderConfigOptions() viper.DecoderConfigOption {
	return func(cfgOpts *mapstructure.DecoderConfig) {
		timeDecoder := func(f reflect.Type,
			t reflect.Type,
			data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(time.Time{}) {
				return data, nil
			}

			return time.Parse(time.Kitchen, data.(string))
		}

		cfgOpts.DecodeHook = mapstructure.ComposeDecodeHookFunc(cfgOpts.DecodeHook, timeDecoder)
	}
}

func (cfg *NotificationCfg) Equal(v *NotificationCfg) bool {
	if cfg.Name != v.Name {
		return false
	}

	if cfg.Message != v.Message {
		return false
	}

	if !reflect.DeepEqual(cfg.CertainTime, v.CertainTime) {
		return false
	}

	if !reflect.DeepEqual(cfg.RandomMessage, v.RandomMessage) {
		return false
	}

	if cfg.RandomTime != v.RandomTime {
		return false
	}

	return true
}
