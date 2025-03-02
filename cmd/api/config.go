package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// getEnv returns environment
func getEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return "dev"
	}
	return env
}

// loadConfig loads configuration for service
func loadConfig(configFileName string) {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
