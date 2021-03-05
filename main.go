package main

import (
	"log"

	"github.com/yakuter/ugin/pkg/config"
	"github.com/yakuter/ugin/pkg/database"
	"github.com/yakuter/ugin/pkg/router"

	_ "github.com/yakuter/ugin/docs"
)

func init() {
	config.Setup()
	database.Setup()
}

// @title GO Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081

func main() {
	config := config.GetConfig()

	db := database.GetDB()
	r := router.Setup(db)

	log.Printf("Server is starting at 127.0.0.1:%s", config.Server.Port)
	r.Run("127.0.0.1:" + config.Server.Port)
}
