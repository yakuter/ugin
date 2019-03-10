package controller

import (
	"fmt"

	"ugin/include"
	"ugin/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

// Post struct alias
type Post = model.Post

func GetPost(c *gin.Context) {
	db = include.GetDB()
	id := c.Params.ByName("id")
	var post Post

	if err := db.Where("id = ? ", id).First(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, post)
	}
}

func GetPosts(c *gin.Context) {
	db = include.GetDB()
	var posts []Post

	if err := db.Find(&posts).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, posts)
	}
}

func CreatePost(c *gin.Context) {
	db = include.GetDB()
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
	db = include.GetDB()
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
	db = include.GetDB()
	id := c.Params.ByName("id")
	var post Post

	if err := db.Where("id = ? ", id).Delete(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, gin.H{"id#" + id: "deleted"})
	}
}
