package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ricejson/rice_chess/internal/repository"
	"github.com/ricejson/rice_chess/internal/repository/dao"
	"github.com/ricejson/rice_chess/internal/service"
	"github.com/ricejson/rice_chess/internal/web"
	"github.com/ricejson/rice_chess/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func main() {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open("root:QWEzxc123456@tcp(localhost:3306)/rice_chess"))
	if err != nil {
		panic(err)
	}
	dao.InitTables(db)
	server := gin.Default()
	// 配置跨域中间件
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "authorization"},
		AllowMethods:     []string{"POST", "OPTIONS"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
	}))
	// 注册路由
	userDAO := dao.NewGORMUserDAO(db)
	userRepository := repository.NewCachedUserRepository(userDAO)
	userService := service.NewUserServiceImpl(userRepository)
	uh := web.NewUserHandler(userService)
	mh := web.NewMatchHandler()
	// 接入登录拦截中间件
	server.Use(middleware.NewLoginMiddlewareBuilder().
		Ignore("/user/login").
		Ignore("/user/register").
		// TODO: 这个路径理应拦截
		Ignore("/match/findMatch").
		Build())
	uh.RegisterRoutes(server)
	mh.RegisterRoutes(server)
	server.Run(":8081")
}
