package service

// func GetPost(db *gorm.DB, id uint) {
// 	var post model.Post
// 	var tags []model.Tag

// 	if err := db.Where("id = ? ", id).First(&post).Error; err != nil {
// 		log.Println(err)
// 		c.AbortWithStatus(404)
// 		return

// 	}

// 	db.Model(&post).Related(&tags)
// 	// SELECT * FROM "tags"  WHERE ("post_id" = 1)

// 	post.Tags = tags
// 	c.JSON(200, post)
// }
