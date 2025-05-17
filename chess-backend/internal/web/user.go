package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/ricejson/rice_chess/internal/service"
	"net/http"
)

const (
	UsernameRegexPattern = `^[\u4e00-\u9fa5a-zA-Z0-9_]{2,8}$`
	PasswordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	UserNameCompile *regexp.Regexp
	PasswordCompile *regexp.Regexp
	userSvc         service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{
		UserNameCompile: regexp.MustCompile(UsernameRegexPattern, regexp.None),
		PasswordCompile: regexp.MustCompile(PasswordRegexPattern, regexp.None),
		userSvc:         userSvc,
	}
}

func (uh *UserHandler) Login(ctx *gin.Context) {
	// 1. 解析请求体参数
	type UserLoginParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var userLoginParams UserLoginParams
	if err := ctx.Bind(&userLoginParams); err != nil {
		return
	}
	// 2. 正则校验格式
	username, pwd := userLoginParams.Username, userLoginParams.Password
	isUsername, err := uh.UserNameCompile.MatchString(username)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isUsername {
		ctx.String(http.StatusOK, "非法用户名格式！")
		return
	}
	isPassword, err := uh.PasswordCompile.MatchString(pwd)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "非法密码格式！")
		return
	}
	// 3. 调用业务层处理登录逻辑
	loginUser, err := uh.userSvc.Login(ctx, username, pwd)
	if err == service.ErrUserNotExists {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户不存在!",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误!",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功!",
		"data": loginUser,
	})
}

func (uh *UserHandler) Register(ctx *gin.Context) {
	// 1. 解析请求体参数
	type UserRegisterParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var userRegisterParams UserRegisterParams
	if err := ctx.Bind(&userRegisterParams); err != nil {
		return
	}
	// 2. 正则校验格式
	username, pwd := userRegisterParams.Username, userRegisterParams.Password
	isUsername, err := uh.UserNameCompile.MatchString(username)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isUsername {
		ctx.String(http.StatusOK, "非法用户名格式！")
		return
	}
	isPassword, err := uh.PasswordCompile.MatchString(pwd)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "非法密码格式！")
		return
	}
	// 3. 调用业务层处理注册逻辑
	err = uh.userSvc.Register(ctx, username, pwd)
	if err == service.ErrUserDuplicate {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户已存在!",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误!",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功!",
	})
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/login", uh.Login)
	ug.POST("/register", uh.Register)
}
