package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HttpHost                  string `mapstructure:"HTTP_HOST"`
	HttpPort                  string `mapstructure:"HTTP_PORT"`
	DatabaseUrl               string `mapstructure:"DATABASE_URL"`
	TonConnect                string `mapstructure:"TON_CONNECT"`
	PlatformSmartContract     string `mapstructure:"PLATFORM_SMART_CONTRACT"`
	SmartContractJettonWallet string `mapstructure:"SMART_CONTRACT_JETTON_WALLET"`
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.HttpHost, c.HttpPort)
}

func (c *Config) ChatAddress() string {
	return net.JoinHostPort(c.HttpHost, c.HttpPort)
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("../.env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	config := new(Config)
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора конфигурации: %v", err)
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	fields := map[string]interface{}{
		"HTTP_HOST":                    config.HttpHost,
		"HTTP_PORT":                    config.HttpPort,
		"DATABASE_URL":                 config.DatabaseUrl,
		"TON_CONNECT":                  config.TonConnect,
		"PLATFORM_SMART_CONTRACT":      config.PlatformSmartContract,
		"SMART_CONTRACT_JETTON_WALLET": config.SmartContractJettonWallet,
	}

	for field, value := range fields {
		if isEmptyValue(value) {
			return fmt.Errorf("отсутствует обязательное поле конфигурации: %s", field)
		}
	}
	return nil
}

func isEmptyValue(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int64:
		return v == 0
	default:
		return false
	}
}
