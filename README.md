# UGin - Ultimate Gin API
UGin is an API boilerplate written in Go (Golang) with Gin Framework. https://github.com/gin-gonic/gin

## Database Support
UGin uses **gorm** as an ORM. So **Sqlite3**, **MySQL** and **PostgreSQL** is supported. You just need to edit **config.yml** file according to your setup. 

**config.yml** content:
```
database:
  driver: "postgres"
  dbname: "ugin"
  username: "user"
  password: "password"
  host: "localhost"
  port: "5432"
```

## Default Models
**UGin** has two models (Post and Tag) as boilerplate to show relational database usage.

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
**UGin** has it's own filtering, search and pagination system. You just need to use these parameters.

**Query parameters:**
```
/posts/?Limit=2
/posts/?Offset=0
/posts/?Sort=ID
/posts/?Order=DESC
/posts/?Search=hello
```

Full: **http://localhost:8081/posts/?Limit=25&Offset=0&Sort=ID&Order=DESC&Search=hello**

## Logging
**UGin** has a very powerful logging logic. There is **application log (ugin.log)**, **database log (ugin.db.log)** and **access log (ugin.access.log)**

### ugin.log:
```
INFO 2021-09-19T00:33:32+03:00 Server is starting at 127.0.0.1:8081
ERROR 2021-09-19T00:39:19+03:00 Failed to open log file ugin.log
```
### ugin.db.log:
```
2021/09/19 00:33:32 /home/user/projects/ugin/pkg/database/database.go:76
[0.023ms] [rows:-] SELECT * FROM `posts` LIMIT 1

2021/09/19 00:33:32 /home/user/go/pkg/mod/gorm.io/driver/sqlite@v1.1.5/migrator.go:261
[0.015ms] [rows:-] SELECT count(*) FROM sqlite_master WHERE type = "index" AND tbl_name = "posts" AND name = "idx_posts_deleted_at"

2021/09/19 00:33:32 /home/user/go/pkg/mod/gorm.io/driver/sqlite@v1.1.5/migrator.go:32
[0.010ms] [rows:-] SELECT count(*) FROM sqlite_master WHERE type='table' AND name="tags"

2021/09/19 00:33:32 /home/user/projects/ugin/pkg/database/database.go:76
[0.011ms] [rows:-] SELECT * FROM `tags` LIMIT 1
```
### ugin.access.log:
```
[GIN] 2021/09/19 - 00:33:43 | 200 |    9.255625ms |       127.0.0.1 | GET      "/posts/"
[GIN] 2021/09/19 - 00:41:51 | 200 |     6.41675ms |       127.0.0.1 | GET      "/posts/4"
```

## Routes
Default **UGin** routes are listed below. 

| METHOD  | ROUTE            | FUNCTION                                                      |
|---------|------------------|---------------------------------------------------------------|
| GET     | /posts/          | github.com/yakuter/ugin/controller.(*Controller).GetPosts     |
| GET     | /posts/:id       | github.com/yakuter/ugin/controller.(*Controller).GetPost      |
| POST    | /posts/          | github.com/yakuter/ugin/controller.(*Controller).CreatePost   |
| PUT     | /posts/:id       | github.com/yakuter/ugin/controller.(*Controller).UpdatePost   |
| DELETE  | /posts/:id       | github.com/yakuter/ugin/controller.(*Controller).DeletePost   |
| GET     | /postsjwt/       | github.com/yakuter/ugin/controller.(*Controller).GetPosts     |
| GET     | /postsjwt/:id    | github.com/yakuter/ugin/controller.(*Controller).GetPost      |
| POST    | /postsjwt/       | github.com/yakuter/ugin/controller.(*Controller).CreatePost   |
| PUT     | /postsjwt/:id    | github.com/yakuter/ugin/controller.(*Controller).UpdatePost   |
| DELETE  | /postsjwt/:id    | github.com/yakuter/ugin/controller.(*Controller).DeletePost   |
| POST    | /auth/signup     | github.com/yakuter/ugin/controller.(*Controller).Signup       |
| POST    | /auth/signin     | github.com/yakuter/ugin/controller.(*Controller).Signin       |
| POST    | /auth/refresh    | github.com/yakuter/ugin/controller.(*Controller).RefreshToken |
| POST    | /auth/check      | github.com/yakuter/ugin/controller.(*Controller).CheckToken   |
| GET     | /admin/dashboard | github.com/yakuter/ugin/controller.Dashboard                  |

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
CORS is important for API's and UGin has it's own CORS middleware in **include/middleware.go**. CORS middleware is called with the code below.
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


## What is next?
- Ugin needs a user service and an authentication method with JWT.