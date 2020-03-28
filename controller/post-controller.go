package controller

import (
	"fmt"
	"log"

	"ugin/model"
	"ugin/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

// Post struct alias
type Post = model.Post

// Tag struct alias
type Tag = model.Tag

// Data is mainle generated for filtering and pagination
type Data struct {
	TotalData    int64
	FilteredData int64
	Data         []Post
}

func GetPost(c *gin.Context) {
	db = database.GetDB()
	id := c.Params.ByName("id")
	var post Post
	var tags []Tag

	if err := db.Where("id = ? ", id).First(&post).Error; err != nil {

		c.AbortWithStatus(404)
		fmt.Println(err)

	} else {

		db.Model(&post).Related(&tags)
		// SELECT * FROM "tags"  WHERE ("post_id" = 1)

		post.Tags = tags
		c.JSON(200, post)
	}
}

func GetPosts(c *gin.Context) {
	db = database.GetDB()
	var posts []Post
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
		c.AbortWithStatus(404)
		log.Println(err)
	} else {
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
}

func CreatePost(c *gin.Context) {
	db = database.GetDB()
	var post Post

	c.BindJSON(&post)

	if err := db.Create(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, post)
	}
}

func UpdatePost(c *gin.Context) {
	db = database.GetDB()
	var post Post
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	c.BindJSON(&post)

	db.Save(&post)
	c.JSON(200, post)
}

func DeletePost(c *gin.Context) {
	db = database.GetDB()
	id := c.Params.ByName("id")
	var post Post

	if err := db.Where("id = ? ", id).Delete(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, gin.H{"id#" + id: "deleted"})
	}
}
