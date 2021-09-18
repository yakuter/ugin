package main

import (
	"github.com/spf13/viper"
	"github.com/yakuter/ugin/pkg/config"
	"github.com/yakuter/ugin/pkg/database"
	"github.com/yakuter/ugin/pkg/logger"
	"github.com/yakuter/ugin/pkg/router"
)

func main() {
	if err := config.Setup(); err != nil {
		logger.Fatalf("config.Setup() error: %s", err)
	}

	if err := database.Setup(); err != nil {
		logger.Fatalf("database.Setup() error: %s", err)
	}

	db := database.GetDB()
	r := router.Setup(db)

	logger.Infof("Server is starting at 127.0.0.1:%s", viper.GetString("server.port"))
	logger.Fatalf("%v", r.Run("127.0.0.1:"+viper.GetString("server.port")))
}
