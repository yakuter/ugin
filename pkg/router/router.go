package router

import (
	"ugin/controller"
	"ugin/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()

	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Non-protected routes
	posts := r.Group("/posts")
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
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"username1": "password1",
		"username2": "password2",
		"username3": "password3",
	}))

	// /admin/dashboard endpoint is now protected
	authorized.GET("/dashboard", controller.Dashboard)

	return r
}
