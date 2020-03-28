package main

import (
	"ugin/pkg/config"
	"ugin/pkg/database"
	"ugin/pkg/router"
)

func init() {
	config.Setup()
	database.Setup()
}

func main() {
	config := config.GetConfig()

	r := router.Setup()
	r.Run("127.0.0.1:" + config.Server.Port)
}
