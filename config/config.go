package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultConfigPath = "/.config/arcane-game-server"
	defaultConfigName = "config.json"
)

func ReadConfigFile(path string) error {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return err

		}

		viper.SetConfigType("json")
		viper.AddConfigPath(home + defaultConfigPath)
		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
