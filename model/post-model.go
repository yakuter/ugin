package model

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Name        string `json:"Name" gorm:"type:varchar(200)"`
	Description string `json:"Description"  gorm:"type:text"`
}

/*
{
    "ID": 1,
    "CreatedAt": "2019-03-10T23:04:15.78791+03:00",
    "UpdatedAt": "2019-03-10T23:04:15.78791+03:00",
    "DeletedAt": null,
    "Name": "Hello World",
    "Description": "This is your first post"
}
*/
