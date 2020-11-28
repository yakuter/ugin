package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/model"
	"github.com/yakuter/ugin/service"
)

var err error

func (base *Controller) GetPost(c *gin.Context) {
	id := c.Params.ByName("id")

	post, err := service.GetPost(base.DB, id)
	if err != nil {
		c.AbortWithStatus(404)
	}

	c.JSON(200, post)
}

func (base *Controller) GetPosts(c *gin.Context) {
	var args model.Args

	// Define and get sorting field
	args.Sort = c.DefaultQuery("Sort", "ID")

	// Define and get sorting order field
	args.Order = c.DefaultQuery("Order", "DESC")

	// Define and get offset for pagination
	args.Offset = c.DefaultQuery("Offset", "0")

	// Define and get limit for pagination
	args.Limit = c.DefaultQuery("Limit", "25")

	// Get search keyword for Search Scope
	args.Search = c.DefaultQuery("Search", "")

	// Fetch results from database
	posts, filteredData, totalData, err := service.GetPosts(c, base.DB, args)
	if err != nil {
		c.AbortWithStatus(404)
	}

	// Fill return data struct
	data := model.Data{
		TotalData:    totalData,
		FilteredData: filteredData,
		Data:         posts,
	}

	c.JSON(200, data)
}

func (base *Controller) CreatePost(c *gin.Context) {
	post := new(model.Post)

	c.ShouldBindJSON(&post)

	post, err := service.SavePost(base.DB, post)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, post)
}

func (base *Controller) UpdatePost(c *gin.Context) {
	id := c.Params.ByName("id")

	post, err := service.GetPost(base.DB, id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.ShouldBindJSON(&post)

	post, err = service.SavePost(base.DB, post)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, post)
}

func (base *Controller) DeletePost(c *gin.Context) {
	id := c.Params.ByName("id")

	err = service.DeletePost(base.DB, id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, gin.H{"id#" + id: "deleted"})
}
