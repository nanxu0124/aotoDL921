package main

import (
	"autoDL921/internal/repository"
	"autoDL921/internal/repository/dao"
	"autoDL921/internal/service"
	"autoDL921/internal/web"
	"autoDL921/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	c := initUser(db)
	c.RegisterRoutes(server)

	server.Run(":8080")
}

func initDB() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:13306)/autoDL921?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://127.0.0.1:3000"},
		//AllowMethods:     []string{},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}
			return strings.Contains(origin, "yumychn@163.com ")
		},
		MaxAge: 12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("autoDL921", store))

	server.Use(middleware.NewSignInMiddlewareBuild().Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	c := web.NewUserHandler(svc)

	return c
}
