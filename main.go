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

	posts := router.Group("/posts")
	{
		posts.GET("/", controller.GetPosts)
		posts.GET("/:id", controller.GetPost)
		posts.POST("/", controller.CreatePost)
		posts.PUT("/:id", controller.UpdatePost)
		posts.DELETE("/:id", controller.DeletePost)
	}

	router.Run(":" + config.Server.Port)
}
