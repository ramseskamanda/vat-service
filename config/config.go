package config

import (
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Configuration struct {
	NumProcessors int
	Port          string
	Environment   string
}

var config *viper.Viper

func Init() {
	config = viper.New()
	config.AutomaticEnv()

	config.SetDefault("PORT", "8080")

	if gin.Mode() == gin.ReleaseMode {
		nuCPU := runtime.NumCPU()
		runtime.GOMAXPROCS(nuCPU)
	}

	config.SetConfigName(gin.Mode())
	config.SetConfigType("yaml")
	config.AddConfigPath("config/")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Error in parsing config file:\n%v", err)
	}
}

func GetString(key string) string {
	return config.GetString(key)
}

func GetBool(key string) bool {
	return config.GetBool(key)
}

func GetConfig() *viper.Viper {
	return config
}
