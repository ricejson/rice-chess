package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ricejson/rice_chess/internal/service"
	"net/http"
	"time"
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
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	if !isUsername {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "非法用户名格式！",
		})
		return
	}
	isPassword, err := uh.PasswordCompile.MatchString(pwd)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	if !isPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "非法密码格式！",
		})
		return
	}
	// 3. 调用业务层处理登录逻辑
	loginUser, err := uh.userSvc.Login(ctx, username, pwd)
	if err == service.ErrUserNotExists {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户不存在！",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	// 4. 登录成功后使用jwt生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)), // 设置过期时间
		},
		Uid: loginUser.UserId,
	})
	tokenStr, err := token.SignedString([]byte("kUD7HXe4bMG6sUWV8pEyQ6JxNQTZkYtu"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	// 存放到header中
	ctx.Header("x-jwt-token", tokenStr)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功！",
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
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	if !isUsername {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "非法用户名格式！",
		})
		return
	}
	isPassword, err := uh.PasswordCompile.MatchString(pwd)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	if !isPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "非法密码格式！",
		})
		return
	}
	// 3. 调用业务层处理注册逻辑
	err = uh.userSvc.Register(ctx, username, pwd)
	if err == service.ErrUserDuplicate {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户已存在！",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功！",
	})
}

func (uh *UserHandler) GetUserProfile(ctx *gin.Context) {
	// 1. 从ctx中取出userId
	value, exists := ctx.Get("userId")
	if !exists {
		// 键不存在
		// 有人在搞你
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	userId := value.(int64)
	// 2. 根据id查找用户
	user, err := uh.userSvc.GetUserInfo(ctx, userId)
	if err == service.ErrUserNotExists {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户不存在！",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误！",
		})
		return
	}
	// 3. 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取用户信息成功！",
		"data": user,
	})
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/login", uh.Login)
	ug.POST("/register", uh.Register)
	ug.GET("/profile", uh.GetUserProfile)
}

// UserClaims 自定义用户载荷
type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
