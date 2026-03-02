package appbuilder

import "github.com/labopase/flevance/packages/utilities"

type AppConfig struct {
	Version     string `mapstructure:"version"`
	Debug       bool   `mapstructure:"debug"`
	Environment string `mapstructure:"environment"`
	Name        string `mapstructure:"name"`
}

func BindConfig[T any](key string) (T, error) {
	return utilities.Bind[T](key, "./", "config", "yaml")
}
