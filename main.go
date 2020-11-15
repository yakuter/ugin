package main

import (
	"log"

	"github.com/yakuter/ugin/pkg/config"
	"github.com/yakuter/ugin/pkg/database"
	"github.com/yakuter/ugin/pkg/router"
)

func init() {
	config.Setup()
	database.Setup()
}

func main() {
	config := config.GetConfig()

	r := router.Setup()
	log.Printf("Server is starting at 127.0.0.1:%s", config.Server.Port)
	r.Run("127.0.0.1:" + config.Server.Port)
}
