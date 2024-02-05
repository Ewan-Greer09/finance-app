package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	API struct {
		ServiceName  string `mapstructure:"service_name"`
		Addr         string `mapstructure:"addr"`
		LogLevel     int    `mapstructure:"log_level"`
		LogFile      string `mapstructure:"log_file"`
		Timeout      int    `mapstructure:"timeout"`
		DatabaseName string `mapstructure:"database_name"`
	} `mapstructure:"api"`
}

func LoadConfig() Config {
	var c Config

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	viper := viper.New()
	viper.AutomaticEnv()

	viper.SetConfigName(os.Getenv("ENV"))
	viper.SetConfigType("json")
	viper.AddConfigPath("./api/config/")

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	return c
}
