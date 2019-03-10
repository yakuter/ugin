# UGin - Ultimate Gin API
UGin is an API boilerplate written in Go (Golang) with Gin Framework. https://www.yakuter.com/

### Dependencies
**UGin** uses **Gin** for main framework, **Gorm** for database and **Viper** for configuration.
```
go get -u github.com/gin-gonic/gin
go get -u github.com/jinzhu/gorm
go get -u github.com/jinzhu/gorm/dialects/postgres
go get -u github.com/jinzhu/gorm/dialects/sqlite
go get -u github.com/jinzhu/gorm/dialects/mysql
go get -u github.com/spf13/viper
```

### Default Endpoints
| Method | URI         | Fonksiyon        |
|--------|-------------|------------------|
| GET    | /posts/     | main.GetPosts    |
| POST   | /posts/     | main.CreatePost  |
| GET    | /posts/:id  | main.GetPost     |
| PUT    | /posts/:id  | main.UpdatePost  |
| DELETE | /posts/:id  | main.DeletePost  |
