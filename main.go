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
	db.AutoMigrate(&model.Post{}, &model.Tag{})

	router := gin.Default()
	router.Use(include.CORS())

	// Non-protected routes
	posts := router.Group("/posts")
	{
		posts.GET("/", controller.GetPosts)
		posts.GET("/:id", controller.GetPost)
		posts.POST("/", controller.CreatePost)
		posts.PUT("/:id", controller.UpdatePost)
		posts.DELETE("/:id", controller.DeletePost)
	}

	// Protected routes
	// For authorized access, group protected routes using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"username1": "password1",
		"username2": "password2",
		"username3": "password3",
	}))

	// /admin/dashboard endpoint is now protected
	authorized.GET("/dashboard", controller.Dashboard)

	router.Run(":" + config.Server.Port)
}
