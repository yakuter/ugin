package controller

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/model"
)

var err error

// Data is mainle generated for filtering and pagination
type Data struct {
	TotalData    int64
	FilteredData int64
	Data         []model.Post
}

func (base *Controller) GetPost(c *gin.Context) {
	db := base.DB
	id := c.Params.ByName("id")
	var post model.Post
	var tags []model.Tag

	if err := db.Where("id = ? ", id).First(&post).Error; err != nil {
		log.Println(err)
		c.AbortWithStatus(404)
		return

	}

	db.Model(&post).Related(&tags)
	// SELECT * FROM "tags"  WHERE ("post_id" = 1)

	post.Tags = tags
	c.JSON(200, post)
}

func (base *Controller) GetPosts(c *gin.Context) {
	db := base.DB
	var posts []model.Post
	var data Data

	// Define and get sorting field
	sort := c.DefaultQuery("Sort", "ID")

	// Define and get sorting order field
	order := c.DefaultQuery("Order", "DESC")

	// Define and get offset for pagination
	offset := c.DefaultQuery("Offset", "0")

	// Define and get limit for pagination
	limit := c.DefaultQuery("Limit", "25")

	// Get search keyword for Search Scope
	search := c.DefaultQuery("Search", "")

	table := "posts"
	query := db.Select(table + ".*")
	query = query.Offset(Offset(offset))
	query = query.Limit(Limit(limit))
	query = query.Order(SortOrder(table, sort, order))
	query = query.Scopes(Search(search))

	if err := query.Find(&posts).Error; err != nil {
		log.Println(err)
		c.AbortWithStatus(404)
		return
	}
	// Count filtered table
	// We are resetting offset to 0 to return total number.
	// This is a fix for Gorm offset issue
	query = query.Offset(0)
	query.Table(table).Count(&data.FilteredData)

	// Count total table
	db.Table(table).Count(&data.TotalData)

	// Set Data result
	data.Data = posts

	c.JSON(200, data)
}

func (base *Controller) CreatePost(c *gin.Context) {
	db := base.DB
	var post model.Post

	c.ShouldBindJSON(&post)

	if err := db.Create(&post).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, post)
}

func (base *Controller) UpdatePost(c *gin.Context) {
	db := base.DB
	var post model.Post
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		log.Println(err)
		c.AbortWithStatus(404)
		return
	}

	c.ShouldBindJSON(&post)

	db.Save(&post)
	c.JSON(200, post)
}

func (base *Controller) DeletePost(c *gin.Context) {
	db := base.DB
	id := c.Params.ByName("id")
	var post model.Post

	if err := db.Where("id = ? ", id).Delete(&post).Error; err != nil {
		log.Println(err)
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, gin.H{"id#" + id: "deleted"})
}
