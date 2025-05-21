package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ricejson/rice_chess/internal/web"
	"net/http"
	"strings"
)

// LoginMiddlewareBuilder 登录中间件
type LoginMiddlewareBuilder struct {
	ignorePaths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Ignore(path string) *LoginMiddlewareBuilder {
	l.ignorePaths = append(l.ignorePaths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 过滤非法请求
		var curPath = ctx.Request.URL.Path
		for _, path := range l.ignorePaths {
			if path == curPath {
				// 不需要登录校验
				return
			}
		}
		// 从请求头中取出token
		var authorization = ctx.GetHeader("authorization")
		if len(authorization) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var splitArr = strings.Split(authorization, " ")
		if len(splitArr) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var tokenStr = splitArr[1]
		// 解析token
		var userClaims web.UserClaims
		parse, err := jwt.ParseWithClaims(tokenStr, &userClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte("kUD7HXe4bMG6sUWV8pEyQ6JxNQTZkYtu"), nil
		})
		if err != nil || !parse.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 将uid放入ctx中
		ctx.Set("userId", userClaims.Uid)
	}
}
