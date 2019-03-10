# UGin - Ultimate Gin API
UGin is an API boilerplate written in Go (Golang) with Gin Framework. https://www.yakuter.com/

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
### Logger and Recovery Middlewares
Gin has 2 important built-in middlewares: Logger and Recovery. UGin calls these two with the followıng code.
```
router := gin.Default()
```

This is same with the following lines.
```
router := gin.New()
router.Use(gin.Logger())
router.Use(gin.Recovery())
```

### CORS Middleware
CORS is really important for API's and UGin has it's own CORS middleware in **include/middleware.go**. CORS middleware is called with the code below.
```
router.Use(include.CORS())
```
There is also a good repo for this: https://github.com/gin-contrib/cors

## Default Endpoints
| Method | URI         | Function         |
|--------|-------------|------------------|
| GET    | /posts/     | main.GetPosts    |
| POST   | /posts/     | main.CreatePost  |
| GET    | /posts/:id  | main.GetPost     |
| PUT    | /posts/:id  | main.UpdatePost  |
| DELETE | /posts/:id  | main.DeletePost  |
