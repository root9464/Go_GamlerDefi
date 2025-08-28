package cfgloader

import (
	"github.com/spf13/viper"
)

type Provider interface {
	Get(key string) any
	Unmarshal(rawVal any) error
	Watch(changes chan<- ConfigChange) error
}

type ConfigChange struct {
	Key   string
	Value any
}

type UniversalConfig struct {
	viper       *viper.Viper
	callbacks   []func(ConfigChange)
	stopWatcher chan struct{}

	watchEnabled bool
}

type Option func(*UniversalConfig)

func WithWatch(status bool) Option {
	return func(uc *UniversalConfig) {
		uc.watchEnabled = status
	}
}
