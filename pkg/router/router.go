package router

import (
	"github.com/jinzhu/gorm"
	"github.com/yakuter/ugin/controller"
	"github.com/yakuter/ugin/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.New()

	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	api := controller.Controller{DB: db}

	// Non-protected routes
	posts := r.Group("/posts")
	{
		posts.GET("/", api.GetPosts)
		posts.GET("/:id", api.GetPost)
		posts.POST("/", api.CreatePost)
		posts.PUT("/:id", api.UpdatePost)
		posts.DELETE("/:id", api.DeletePost)
	}

	// JWT-protected routes
	postsjwt := r.Group("/postsjwt", middleware.Authorize())
	{
		postsjwt.GET("/", api.GetPosts)
		postsjwt.GET("/:id", api.GetPost)
		postsjwt.POST("/", api.CreatePost)
		postsjwt.PUT("/:id", api.UpdatePost)
		postsjwt.DELETE("/:id", api.DeletePost)
	}

	authRouter := r.Group("/auth")
	{
		authRouter.POST("/signup", api.Signup)
		authRouter.POST("/signin", api.Signin)
		authRouter.POST("/refresh", api.RefreshToken)
		authRouter.POST("/check", api.CheckToken)
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
