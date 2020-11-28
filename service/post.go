package service

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yakuter/ugin/model"
)

func GetPost(db *gorm.DB, id string) (*model.Post, error) {
	var err error
	post := new(model.Post)

	if err := db.Where("id = ? ", id).Preload("Tags").First(&post).Error; err != nil {
		log.Println(err)

		return nil, err
	}

	return post, err
}

func GetPosts(c *gin.Context, db *gorm.DB, args model.Args) ([]model.Post, int64, int64, error) {
	posts := []model.Post{}
	var filteredData, totalData int64

	table := "posts"
	query := db.Select(table + ".*")
	query = query.Offset(Offset(args.Offset))
	query = query.Limit(Limit(args.Limit))
	query = query.Order(SortOrder(table, args.Sort, args.Order))
	query = query.Scopes(Search(args.Search))

	if err := query.Preload("Tags").Find(&posts).Error; err != nil {
		log.Println(err)
		return posts, filteredData, totalData, err
	}

	// // Count filtered table
	// // We are resetting offset to 0 to return total number.
	// // This is a fix for Gorm offset issue
	query = query.Offset(0)
	query.Table(table).Count(&filteredData)

	// // Count total table
	db.Table(table).Count(&totalData)

	return posts, filteredData, totalData, nil
}

// SavePost both cretes and updates post according to if ID field is empty or not
func SavePost(db *gorm.DB, post *model.Post) (*model.Post, error) {
	if err := db.Save(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}

// DeletePost soft deletes all records.
func DeletePost(db *gorm.DB, id string) error {
	post := new(model.Post)
	if err := db.Where("id = ? ", id).Delete(&post).Error; err != nil {
		log.Println(err)
		return err
	}

	tag := new(model.Tag)
	if err := db.Where("post_id = ? ", id).Delete(&tag).Error; err != nil {
		log.Println(err)
	}

	return nil
}
