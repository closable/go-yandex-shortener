package config

import (
	"github.com/caarlos0/env/v10"
)

type config struct {
	//ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	ServerAddress string `env:"SERVER_ADDRESS"`
	//BaseURL       string `env:"BASE_URL" envDefault:"localhost:8080"`
	BaseURL string `env:"BASE_URL"`
}

var ConfigEnv = config{}

func ParseConfigEnv() {

	env.Parse(&ConfigEnv)

}

func GetEnvParam(typeVar string) string {

	if typeVar == "RUN_SERVER" {
		if len(ConfigEnv.ServerAddress) > 0 {
			return ConfigEnv.ServerAddress
		}
		if len(FlagRunAddr) > 0 {
			return FlagRunAddr
		}
	}

	if typeVar == "SND_SERVER" {
		if len(ConfigEnv.BaseURL) > 0 {
			return ConfigEnv.BaseURL
		}
		if len(FlagSendAddr) > 0 {
			return FlagSendAddr
		}
	}
	return "localhost:8080"
}
