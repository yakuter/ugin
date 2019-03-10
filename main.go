package main

import (
	"ugin/config"
	"ugin/controller"
	"ugin/include"
	"ugin/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func main() {
	config := config.InitConfig()

	db = include.InitDB()
	defer db.Close()
	db.AutoMigrate(&model.Post{})

	router := gin.Default()
	router.Use(include.CORS())

	v1 := router.Group("/v1")
	{
		v1.GET("/", controller.GetPosts)
		v1.GET("/:id", controller.GetPost)
		v1.POST("/", controller.CreatePost)
		v1.PUT("/:id", controller.UpdatePost)
		v1.DELETE("/:id", controller.DeletePost)
	}

	router.Run(":" + config.Server.Port)
}
