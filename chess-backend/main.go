package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ricejson/rice_chess/internal/repository"
	"github.com/ricejson/rice_chess/internal/repository/dao"
	"github.com/ricejson/rice_chess/internal/service"
	"github.com/ricejson/rice_chess/internal/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open("root:QWEzxc123456@tcp(localhost:3306)/rice_chess"))
	if err != nil {
		panic(err)
	}
	dao.InitTables(db)
	server := gin.Default()
	// 注册路由
	userDAO := dao.NewGORMUserDAO(db)
	userRepository := repository.NewCachedUserRepository(userDAO)
	userService := service.NewUserServiceImpl(userRepository)
	uh := web.NewUserHandler(userService)
	uh.RegisterRoutes(server)
	server.Run(":8081")
}
