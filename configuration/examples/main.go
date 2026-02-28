package main

import (
	"fmt"

	"github.com/labopase/flevance/configuration"
)

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type App struct {
	Debug   bool   `mapstructure:"debug"`
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

type Config struct {
	App    App    `mapstructure:"app"`
	Server Server `mapstructure:"server"`
}

func main() {
	config := BindAllConfig()
	fmt.Println(config)

	appConfig := BindAppConfig()
	fmt.Println(appConfig)

	serverConfig := BindServerConfig()
	fmt.Println(serverConfig)
}

func BindAllConfig() Config {
	config, err := configuration.Bind[Config]("", "./", "config", "yaml")
	if err != nil {
		panic(err)
	}

	return config
}

func BindServerConfig() Server {
	config, err := configuration.Bind[Server]("server", "./", "config", "yaml")
	if err != nil {
		panic(err)
	}

	return config
}

func BindAppConfig() App {
	config, err := configuration.Bind[App]("app", "./", "config", "yaml")
	if err != nil {
		panic(err)
	}

	return config
}
