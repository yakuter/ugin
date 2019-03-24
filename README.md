# UGin - Ultimate Gin API
UGin is an API boilerplate written in Go (Golang) with Gin Framework. https://github.com/gin-gonic/gin

## Database Support
UGin uses **gorm** as an ORM. So **Sqlite3**, **MySQL** and **PostgreSQL** is supported. You just need to edit **config.yml** file according to your setup. 

**config.yml** content:
```
database:
  driver: "postgres"
  dbname: "database_name"
  username: "user"
  password: "password"
  host: "localhost"
  port: "5432"
```

## Default Models
UGin has two models (Post and Tag) as boilerplate to show relational database usage.

**/model/post-model.go** content:
```
type Post struct {
	gorm.Model
	Name        string `json:"Name" gorm:"type:varchar(255)"`
	Description string `json:"Description"  gorm:"type:text"`
	Tags        []Tag  // One-To-Many relationship (has many - use Tag's UserID as foreign key)
}

type Tag struct {
	gorm.Model
	PostID      uint   `gorm:"index"` // Foreign key (belongs to)
	Name        string `json:"Name" gorm:"type:varchar(255)"`
	Description string `json:"Description" gorm:"type:text"`
}
```

## Filtering, Search and Pagination
UGin has it's own filtering, search and pagination system. You just need to use these parameters.

**Query parameters:**
```
/posts/?limit=2
/posts/?offset=0
/posts/?name=Third
/posts/?description=My
/posts/?order=name|asc
```

Full: **http://localhost:8081/posts/?limit=2&offset=0&name=Third&description=My&order=name|asc**

## Dependencies
**UGin** uses **Gin** for main framework, **Gorm** for database and **Viper** for configuration.
```
go get -u github.com/gin-gonic/gin
go get -u github.com/jinzhu/gorm
go get -u github.com/jinzhu/gorm/dialects/postgres
go get -u github.com/jinzhu/gorm/dialects/sqlite
go get -u github.com/jinzhu/gorm/dialects/mysql
go get -u github.com/spf13/viper
```
## Middlewares
### 1. Logger and Recovery Middlewares
Gin has 2 important built-in middlewares: **Logger** and **Recovery**. UGin calls these two in default.
```
router := gin.Default()
```

This is same with the following lines.
```
router := gin.New()
router.Use(gin.Logger())
router.Use(gin.Recovery())
```

### 2. CORS Middleware
CORS is really important for API's and UGin has it's own CORS middleware in **include/middleware.go**. CORS middleware is called with the code below.
```
router.Use(include.CORS())
```
There is also a good repo for this: https://github.com/gin-contrib/cors

### 3. BasicAuth Middleware
Almost every API needs a protected area. Gin has **BasicAuth** middleware for protecting routes. Basic Auth is an authorization type that requires a verified username and password to access a data resource. In UGin, you can find an example for a basic auth. To access these protected routes, you need to add **Basic Authorization credentials** in your requests. If you try to reach these endpoints from browser, you should see a window prompting you for username and password.

```
authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
    "username": "password",
}))

// /admin/dashboard endpoint is now protected
authorized.GET("/dashboard", controller.Dashboard)
```

If you want to use JWT for authorization in UGin, you can check this: https://github.com/appleboy/gin-jwt

## Default Endpoints
| Method | URI              | Function               |
|--------|------------------|------------------------|
| GET    | /posts/          | controller.GetPosts    |
| POST   | /posts/          | controller.CreatePost  |
| GET    | /posts/:id       | controller.GetPost     |
| PUT    | /posts/:id       | controller.UpdatePost  |
| DELETE | /posts/:id       | controller.DeletePost  |
| GET    | /admin/dashboard | controller.Dashboard   |

## Demonstration and Documentation
I used [gencebay](https://github.com/gencebay)'s great API request and response mock tool [httplive](https://github.com/gencebay/httplive) to show the endpoints and responses of UGin.  So there is a **httplive** folder which is including **httplive.db**. This folder is not necessary for API, you can delete it whenever you want. Just wanted to give you **httplive.db** that can be used in your local httplive installation.

![httplive](https://github.com/yakuter/ugin/blob/461e29c340471acaccb6bd8f6b988939eadadee1/httplive/httplive-screenshot.png)
